package handle

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dopamine-joker/zu_logic/db"
	"github.com/dopamine-joker/zu_logic/handle/dao"
	"github.com/dopamine-joker/zu_logic/misc"
	"github.com/dopamine-joker/zu_logic/proto"
	"github.com/fatih/structs"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type RpcLogicServer struct {
	proto.UnimplementedRpcLogicServiceServer

	Addr string
}

func NewRpcLogicServer(Host string) *RpcLogicServer {
	return &RpcLogicServer{
		Addr: Host,
	}
}

//Login 非token登陆
func (r *RpcLogicServer) Login(ctx context.Context, request *proto.LoginRequest) (*proto.LoginResponse, error) {
	response := &proto.LoginResponse{
		Code: misc.CodeFail,
	}
	email := request.GetEmail()
	// 查找数据库是否存在该用户
	user := dao.GetUserByEmail(ctx, request.GetEmail())
	if user.Id == 0 {
		misc.Logger.Error("login err, user does not exist", zap.String("email", email))
		return response, errors.New("login err, user does not exist")
	}

	// 密码是否正确
	if request.GetPassword() != user.Password {
		return response, errors.New("password err")
	}

	// redis已有相关tokenkey,则删除
	prefix := misc.GetTokenKeyPrefix(user.Id)
	var cursor uint64
	for {
		var keys []string
		var err error
		keys, cursor, err = db.RedisClient.Scan(ctx, cursor, fmt.Sprintf("%s*", prefix), 20).Result()
		if err != nil {
			return response, err
		}
		for _, key := range keys {
			misc.Logger.Info("scan key", zap.String("key", key))
			if _, err = db.RedisClient.Del(ctx, key).Result(); err != nil {
				misc.Logger.Warn("del key err", zap.String("key", key))
			}
		}
		if cursor == 0 {
			break
		}
	}

	// 生成新的token
	tokenId := misc.GenerateTokenKey(user.Id)

	// redis写入token
	if _, err := db.RedisClient.Pipelined(ctx, func(pipeliner redis.Pipeliner) error {
		dataMap := structs.Map(user)
		// 保存token
		pipeliner.HSet(ctx, tokenId, dataMap)
		pipeliner.Expire(ctx, tokenId, 86400*time.Second)
		// 写入token
		//pipeliner.Set(ctx, loginTokenKey, tokenId, 86400*time.Second)
		return nil
	}); err != nil {
		misc.Logger.Warn("set token err", zap.Int32("userId", user.Id), zap.Error(err))
		return response, err
	}

	response.Code = misc.CodeSuccess
	response.AuthToken = tokenId
	response.User = &proto.User{
		Id:     user.Id,
		Email:  user.Email,
		Phone:  user.Phone,
		Name:   user.Name,
		Face:   user.Face,
		School: user.School,
		Sex:    user.Sex,
	}
	return response, nil
}

//TokenLogin 使用token登陆
func (r *RpcLogicServer) TokenLogin(ctx context.Context, request *proto.TokenLoginRequest) (*proto.TokenLoginResponse, error) {

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("token", request.GetToken()))

	response := &proto.TokenLoginResponse{
		Code: misc.CodeFail,
	}
	tokenId := request.GetToken()
	// 解析token是否有效
	var err error
	var num int64
	// 查询redis的token是否存在(即过期，或根本就不存在)
	if num, err = db.RedisClient.Exists(ctx, tokenId).Result(); err != nil {
		misc.Logger.Error("token login redis exists err", zap.String("token", tokenId), zap.Error(err))
		return response, err
	}
	// num=0说明token不存在
	if num == 0 {
		return response, errors.New("token not exists")
	}

	// 刷新token时间
	if _, err = db.RedisClient.Expire(ctx, tokenId, 86400*time.Second).Result(); err != nil {
		misc.Logger.Error("redis expire token key err", zap.String("tokenId", tokenId), zap.Error(err))
		return response, err
	}

	userDataMap, err := db.RedisClient.HGetAll(ctx, tokenId).Result()
	if err != nil {
		misc.Logger.Error("redis hgetall err", zap.Error(err))
		return response, err
	}
	user := &proto.User{}
	intUserId, _ := strconv.Atoi(userDataMap["Id"])
	user.Id = int32(intUserId)
	user.Email = userDataMap["Email"]
	user.Name = userDataMap["Name"]
	user.Face = userDataMap["Face"]
	user.Phone = userDataMap["Phone"]
	user.School = userDataMap["School"]
	sex, err := strconv.ParseInt(userDataMap["Sex"], 10, 32)
	if err != nil {
		misc.Logger.Error("redis get user info err", zap.Error(err))
		return response, err
	}
	user.Sex = int32(sex)

	response.Code = misc.CodeSuccess
	response.AuthToken = tokenId
	response.User = user
	return response, nil
}

func (r *RpcLogicServer) UpdateUser(ctx context.Context, req *proto.UpdateUserRequest) (*proto.UpdateUserResponse, error) {
	response := &proto.UpdateUserResponse{
		Code: misc.CodeFail,
	}
	if err := dao.UpdateUser(ctx, req.Email, req.Phone, req.Name, req.Password, req.School, req.Sex, req.Uid); err != nil {
		misc.Logger.Error("update user err", zap.Error(err))
		return response, err
	}

	response.Code = misc.CodeSuccess
	return response, nil
}

//Register 注册
func (r *RpcLogicServer) Register(ctx context.Context, request *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	response := &proto.RegisterResponse{
		Code: misc.CodeFail,
	}
	user := dao.GetUserByEmail(ctx, request.GetEmail())
	if user.Id > 0 {
		misc.Logger.Warn("register err, user have already exist", zap.String("email", request.GetEmail()))
		return response, errors.New("email have been register, please login")
	}
	// 数据库增加一个user
	userId, err := dao.AddUser(ctx, request.GetEmail(), request.GetName(), request.GetPassword())
	if err != nil {
		misc.Logger.Error("add user err", zap.String("email", request.GetEmail()), zap.Error(err))
		return response, err
	}
	if userId == 0 {
		misc.Logger.Error("register userId empty", zap.String("user name", request.GetName()))
		return response, errors.New("register userId empty")
	}

	response.Code = misc.CodeSuccess
	return response, nil
}

//CheckAuth 检查token
func (r *RpcLogicServer) CheckAuth(ctx context.Context, request *proto.CheckAuthRequest) (*proto.CheckAuthResponse, error) {
	response := &proto.CheckAuthResponse{
		Code: misc.CodeFail,
	}
	tokenId := request.GetAuthToken()
	userDataMap, err := db.RedisClient.HGetAll(ctx, tokenId).Result()
	if err != nil {
		misc.Logger.Error("check auth fail", zap.String("token", tokenId))
		return response, err
	}

	user := &proto.User{}
	intUserId, _ := strconv.Atoi(userDataMap["Id"])
	user.Id = int32(intUserId)
	user.Email = userDataMap["Email"]
	user.Name = userDataMap["Name"]

	response.Code = misc.CodeSuccess
	response.AuthToken = tokenId
	response.User = user
	return response, nil
}

func (r *RpcLogicServer) Logout(ctx context.Context, request *proto.LogoutRequest) (*proto.LogoutResponse, error) {
	response := &proto.LogoutResponse{
		Code: misc.CodeFail,
	}
	tokenId := request.GetToken()
	// 解析token是否有效
	var err error
	var num int64
	// 查询redis的token是否存在(即过期，或根本就不存在)
	if num, err = db.RedisClient.Exists(ctx, tokenId).Result(); err != nil {
		misc.Logger.Error("logout redis exists err", zap.String("token", tokenId), zap.Error(err))
		return response, err
	}
	// num=0说明token不存在
	if num == 0 {
		return response, errors.New("token not exists")
	}

	if _, err = db.RedisClient.Del(ctx, request.GetToken()).Result(); err != nil {
		misc.Logger.Error("logout redis del err", zap.String("token", tokenId), zap.Error(err))
		return response, err
	}

	response.Code = misc.CodeSuccess
	return response, nil
}

const (
	fileNamePrefix  = "pic_"
	coverNamePrefix = "cover_"
	facePrefix      = "face_"
	voicePrefix     = "voice_"
	tmpFilePrePath  = "./upload/"
)

//writeFile 根据路径保存图片到本地
func writeFile(ctx context.Context, file *proto.FileStream, namePrefix string) (key string, path string, err error) {
	//得到图片的类型,png,jpg,jpeg等
	fileType := strings.Split(file.GetName(), ".")[1]
	name := fmt.Sprintf("%s.%s", uuid.New().String(), fileType)
	//根据路径名创建一个临时文件
	tmpFile, err := os.CreateTemp(tmpFilePrePath, fmt.Sprintf("%s*_%s", namePrefix, name))
	defer func() {
		_ = tmpFile.Close()
		_ = os.Remove(tmpFile.Name())
	}()
	if err != nil {
		return "", "", err
	}
	//内容写到临时文件里
	if _, err = tmpFile.Write(file.Content); err != nil {
		return "", "", err
	}

	//临时文件上传cos
	res, _, err := misc.CosClient.Object.Upload(ctx, name, tmpFile.Name(), nil)
	if err != nil {
		misc.Logger.Error("cos upload err", zap.Error(err))
		return "", "", err
	}

	return name, res.Location, nil
}

func (r *RpcLogicServer) UploadFace(ctx context.Context, req *proto.UploadFaceRequest) (*proto.UploadFaceResponse, error) {
	response := &proto.UploadFaceResponse{
		Code: misc.CodeFail,
	}

	_, facePath, err := writeFile(ctx, req.Pic, facePrefix)
	if err != nil {
		return response, err
	}

	if err = dao.UpdateFace(ctx, facePath, req.Uid); err != nil {
		return response, err
	}

	log.Println(facePath)

	response.Code = misc.CodeSuccess
	response.Path = facePath
	return response, nil
}

func (r *RpcLogicServer) VoiceToTxt(ctx context.Context, req *proto.VoiceToTxtRequest) (*proto.VoiceToTxtResponse, error) {
	response := &proto.VoiceToTxtResponse{
		Code: misc.CodeFail,
	}
	// 写cos
	key, path, err := writeFile(ctx, req.VoiceFile, voicePrefix)
	if err != nil {
		return response, err
	}
	// 调sdk包接口
	txt, err := processVoice(path)
	if err != nil {
		return response, err
	}
	// 删除cos相关文件
	removeCosVoice(ctx, key)
	// 封装文字
	response.Code = misc.CodeSuccess
	response.Txt = txt
	return response, nil
}

func removeCosVoice(ctx context.Context, key string) {
	_, err := misc.CosClient.Object.Delete(ctx, key)
	if err != nil {
		misc.Logger.Error("cos delete err", zap.Error(err))
	}
}

//UploadPic 上传图片到logic
func (r *RpcLogicServer) UploadPic(ctx context.Context, req *proto.UploadRequest) (*proto.UploadResponse, error) {
	response := &proto.UploadResponse{
		Code: misc.CodeFail,
	}

	//写图片到cos，并把url保存起来
	var fileUrlList []string
	for _, file := range req.PicList {
		_, path, err := writeFile(ctx, file, fileNamePrefix)
		if err != nil {
			return response, err
		}
		fileUrlList = append(fileUrlList, path)
	}

	//提取封面图片
	_, coverPath, err := writeFile(ctx, req.Cover, coverNamePrefix)
	if err != nil {
		return response, err
	}

	log.Println(coverPath)

	uid := req.Uid
	name := req.Name
	price, err := strconv.ParseFloat(req.Price, 64)
	if err != nil {
		return response, errors.New("price can not be parse")
	}
	detail := req.Detail

	//数据写数据库,包括物品信息,图片等
	goodsId, err := dao.AddGoods(ctx, name, detail, price, uid, coverPath, fileUrlList)
	if err != nil {
		return response, err
	}

	misc.Logger.Info("add goods into sql success", zap.Int32("goodsId", goodsId))

	response.Code = misc.CodeSuccess
	response.GoodId = goodsId
	return response, nil
}

func (r *RpcLogicServer) DeleteGoods(ctx context.Context, req *proto.DeleteGoodsRequest) (*proto.DeleteGoodsResponse, error) {
	response := &proto.DeleteGoodsResponse{
		Code: misc.CodeFail,
	}
	err := dao.DelGoods(ctx, req.Gid)
	if err != nil {
		return response, err
	}
	response.Code = misc.CodeSuccess
	return response, nil
}

func (r *RpcLogicServer) GetGoods(ctx context.Context, req *proto.GetGoodsRequest) (*proto.GetGoodsResponse, error) {
	response := &proto.GetGoodsResponse{
		Code: misc.CodeFail,
	}

	page := req.GetPage()
	count := req.GetCount()

	goodsList, err := dao.GetGoods(ctx, page, count)
	if err != nil {
		return response, err
	}

	var protoList []*proto.Goods

	for _, goods := range goodsList {
		//增加到请求的商品列表
		protoList = append(protoList, &proto.Goods{
			Id:      goods.Id,
			Name:    goods.Name,
			Uname:   goods.Uname,
			Price:   strconv.FormatFloat(goods.Price, 'f', 2, 32),
			SellNum: goods.SellNum,
			Cover:   goods.Cover,
		})
	}

	response.GoodsList = protoList
	response.Code = misc.CodeSuccess
	return response, nil
}

func (r *RpcLogicServer) UserGoods(ctx context.Context, req *proto.GetUserGoodsListRequest) (*proto.GetUserGoodsListResponse, error) {
	response := &proto.GetUserGoodsListResponse{
		Code: misc.CodeFail,
	}

	list, err := dao.GetGoodsByUserId(ctx, req.Uid)
	if err != nil {
		return response, err
	}

	var protoList []*proto.GoodsDetail

	for _, l := range list {
		protoList = append(protoList, &proto.GoodsDetail{
			Gid:        l.Id,
			Uid:        l.UserId,
			Name:       l.Name,
			Uname:      l.Uname,
			Price:      strconv.FormatFloat(l.Price, 'f', 2, 32),
			SellNum:    l.SellNum,
			Detail:     l.Detail,
			Cover:      l.Cover,
			CreateTime: l.CreateTime.Unix(),
		})
	}

	response.Code = misc.CodeSuccess
	response.List = protoList
	return response, nil
}

func (r *RpcLogicServer) GetGoodsPic(ctx context.Context, req *proto.GetGoodsDetailRequest) (*proto.GetGoodsDetailResponse, error) {
	response := &proto.GetGoodsDetailResponse{
		Code: misc.CodeFail,
	}
	goods, err := dao.GetGoodsDetail(ctx, req.GetGid())
	if err != nil {
		misc.Logger.Error("get goods pic err", zap.Error(err))
		return response, err
	}

	var picList []*proto.Pic

	for _, p := range goods.PicList {
		picList = append(picList, &proto.Pic{
			Pid:  p.Id,
			Path: p.Path,
		})
	}

	response.Code = misc.CodeSuccess
	response.Goods = &proto.GoodsDetail{
		Gid:        goods.Id,
		Uid:        goods.UserId,
		Name:       goods.Name,
		Uname:      goods.Uname,
		Price:      strconv.FormatFloat(goods.Price, 'f', 2, 32),
		SellNum:    goods.SellNum,
		Detail:     goods.Detail,
		Cover:      goods.Cover,
		CreateTime: goods.CreateTime.Unix(),
	}
	response.PicList = picList

	return response, nil
}

func (r *RpcLogicServer) SearchGoods(ctx context.Context, req *proto.SearchGoodsRequest) (*proto.SearchGoodsResponse, error) {
	response := &proto.SearchGoodsResponse{
		Code: misc.CodeFail,
	}

	list, err := dao.GetGoodsByName(ctx, req.Name)
	if err != nil {
		misc.Logger.Error("get goods by name err", zap.Error(err))
		return response, err
	}

	var protoList []*proto.GoodsDetail

	for _, l := range list {
		protoList = append(protoList, &proto.GoodsDetail{
			Gid:        l.Id,
			Uid:        l.UserId,
			Name:       l.Name,
			Uname:      l.Uname,
			Price:      strconv.FormatFloat(l.Price, 'f', 2, 32),
			SellNum:    l.SellNum,
			Detail:     l.Detail,
			Cover:      l.Cover,
			CreateTime: l.CreateTime.Unix(),
		})
	}

	response.List = protoList
	response.Code = misc.CodeSuccess

	return response, nil
}

func (r *RpcLogicServer) AddOrder(ctx context.Context, req *proto.AddOrderRequest) (*proto.AddOrderResponse, error) {
	response := &proto.AddOrderResponse{
		Code: misc.CodeFail,
	}

	// 往redis添加记录
	redisOrder := dao.RedisOrderAdd{
		BuyId:  req.Buyid,
		SellId: req.Sellid,
		GId:    req.Gid,
		School: req.School,
		Status: dao.COMMIT,
	}
	addOrderMsg, err := json.Marshal(redisOrder)
	if err != nil {
		return response, err
	}
	if err := db.RedisClient.LPush(ctx, db.RedisOrderAdd, addOrderMsg).Err(); err != nil {
		misc.Logger.Error("add order err", zap.Error(err))
		return response, err
	}

	//oid, err := dao.AddOrder(ctx, req.Buyid, req.Sellid, req.Gid, dao.COMMIT)
	//if err != nil {
	//	misc.Logger.Error("add order err", zap.Error(err))
	//	return response, err
	//}

	response.Code = misc.CodeSuccess
	return response, nil
}

func (r *RpcLogicServer) GetBuyOrder(ctx context.Context, req *proto.GetBuyOrderRequest) (*proto.GetBuyOrderResponse, error) {
	response := &proto.GetBuyOrderResponse{
		Code: misc.CodeFail,
	}

	list, err := dao.GetBuyOrder(ctx, req.Buyid)
	if err != nil {
		misc.Logger.Error("add order err", zap.Error(err))
		return response, err
	}

	var protoList []*proto.Order
	for _, o := range list {
		protoList = append(protoList, &proto.Order{
			Id:       o.Id,
			Buyid:    o.BuyId,
			BuyName:  o.BuyUserName,
			Sellid:   o.SellId,
			SellName: o.SellUserName,
			GId:      o.GId,
			Gname:    o.GName,
			School:   o.School,
			Price:    strconv.FormatFloat(o.Price, 'f', 2, 32),
			Cover:    o.Cover,
			Status:   o.Status,
			Time:     o.Time.Unix(),
		})
	}

	response.Code = misc.CodeSuccess
	response.OrderList = protoList
	return response, nil
}

func (r *RpcLogicServer) GetSellOrder(ctx context.Context, req *proto.GetSellOrderRequest) (*proto.GetSellOrderResponse, error) {
	response := &proto.GetSellOrderResponse{
		Code: misc.CodeFail,
	}

	list, err := dao.GetSellOrder(ctx, req.Sellid)
	if err != nil {
		misc.Logger.Error("add order err", zap.Error(err))
		return response, err
	}

	var protoList []*proto.Order
	for _, o := range list {
		protoList = append(protoList, &proto.Order{
			Id:       o.Id,
			Buyid:    o.BuyId,
			BuyName:  o.BuyUserName,
			Sellid:   o.SellId,
			SellName: o.SellUserName,
			GId:      o.GId,
			Gname:    o.GName,
			School:   o.School,
			Price:    strconv.FormatFloat(o.Price, 'f', 2, 32),
			Cover:    o.Cover,
			Status:   o.Status,
			Time:     o.Time.Unix(),
		})
	}

	response.Code = misc.CodeSuccess
	response.OrderList = protoList
	return response, nil
}

func (r *RpcLogicServer) UpdateOrder(ctx context.Context, req *proto.UpdateOrderRequest) (*proto.UpdateOrderResponse, error) {
	response := &proto.UpdateOrderResponse{
		Code: misc.CodeFail,
	}

	status := dao.OrderStatus(req.Status)

	redisOrder := dao.RedisOrderUpdate{
		OrderId: req.Id,
		Status:  status,
	}
	orderUpdateMsg, err := json.Marshal(redisOrder)
	if err != nil {
		return response, err
	}
	if err := db.RedisClient.LPush(ctx, db.RedisOrderUpdate, orderUpdateMsg).Err(); err != nil {
		misc.Logger.Error("add order err", zap.Error(err))
		return response, err
	}

	misc.Logger.Info("put update order to redis", zap.Int32("orderId", req.Id))

	//if err := dao.UpdateOrder(ctx, req.Id, status); err != nil {
	//	misc.Logger.Error("update order err", zap.Error(err))
	//	return response, err
	//}

	response.Code = misc.CodeSuccess
	return response, nil
}

func (r *RpcLogicServer) AddFavorites(ctx context.Context, req *proto.AddFavoritesRequest) (*proto.AddFavoritesResponse, error) {
	response := &proto.AddFavoritesResponse{
		Code: misc.CodeFail,
	}

	fid, err := dao.AddFavorites(ctx, req.Uid, req.Gid)
	if err != nil {
		misc.Logger.Error("add favorites err", zap.Error(err))
		return response, err
	}

	response.Code = misc.CodeSuccess
	response.Fid = fid
	return response, nil
}

func (r *RpcLogicServer) DeleteFavorites(ctx context.Context, req *proto.DeleteFavoritesRequest) (*proto.DeleteFavoritesResponse, error) {
	response := &proto.DeleteFavoritesResponse{
		Code: misc.CodeFail,
	}
	if err := dao.DeleteFavorites(ctx, req.Fid); err != nil {
		return response, err
	}
	response.Code = misc.CodeSuccess
	return response, nil
}

func (r *RpcLogicServer) GetUserFavorites(ctx context.Context, req *proto.GetUserFavoritesRequest) (*proto.GetUserFavoritesResponse, error) {
	response := &proto.GetUserFavoritesResponse{
		Code: misc.CodeFail,
	}
	list, err := dao.GetUserFavorites(ctx, req.Uid)
	if err != nil {
		misc.Logger.Error("get user favorites err", zap.Error(err))
		return response, err
	}
	var favoritesList []*proto.UserFavorites
	for _, favorites := range list {
		favoritesList = append(favoritesList, &proto.UserFavorites{
			Id:    favorites.Id,
			Uid:   favorites.UId,
			Gid:   favorites.GId,
			Name:  favorites.Name,
			Price: strconv.FormatFloat(favorites.Price, 'f', 2, 32),
			Cover: favorites.Cover,
		})
	}
	response.UserFavoritesList = favoritesList
	response.Code = misc.CodeSuccess
	return response, nil
}

func (r *RpcLogicServer) AddComment(ctx context.Context, req *proto.AddCommentRequest) (*proto.AddCommentResponse, error) {
	response := &proto.AddCommentResponse{
		Code: misc.CodeFail,
	}
	cid, err := dao.AddComment(ctx, req.Uid, req.Gid, req.Oid, req.Level, req.Content)
	if err != nil {
		misc.Logger.Error("add comment err", zap.Error(err))
		return response, err
	}
	response.Code = misc.CodeSuccess
	response.Cid = cid
	return response, nil
}

func (r *RpcLogicServer) GetCommentByUserId(ctx context.Context, req *proto.GetCommentByUserIdRequest) (*proto.GetCommentByUserIdResponse, error) {
	response := &proto.GetCommentByUserIdResponse{
		Code: misc.CodeFail,
	}
	list, err := dao.GetCommentByUserId(ctx, req.Uid)
	if err != nil {
		misc.Logger.Error("get comment list by userId err", zap.Error(err))
		return response, err
	}
	var userCommentList []*proto.UserComment
	for _, comment := range list {
		userCommentList = append(userCommentList, &proto.UserComment{
			Id:      comment.Id,
			Uid:     comment.UId,
			Gid:     comment.GId,
			Oid:     comment.Oid,
			Content: comment.Content,
			Level:   comment.Level,
			Time:    comment.TIme.Unix(),
			Name:    comment.Name,
			Price:   strconv.FormatFloat(comment.Price, 'f', 2, 32),
			Cover:   comment.Cover,
		})
	}
	response.Code = misc.CodeSuccess
	response.UserCommentList = userCommentList
	return response, nil
}

func (r *RpcLogicServer) GetCommentByGoodsId(ctx context.Context, req *proto.GetCommentByGoodsIdRequest) (*proto.GetCommentByGoodsIdResponse, error) {
	response := &proto.GetCommentByGoodsIdResponse{
		Code: misc.CodeFail,
	}
	list, err := dao.GetCommentByGoodsId(ctx, req.Gid)
	if err != nil {
		misc.Logger.Error("get comment by goods id err", zap.Error(err))
		return response, err
	}
	var commentList []*proto.GoodsComment
	for _, comment := range list {
		commentList = append(commentList, &proto.GoodsComment{
			Id:      comment.Id,
			Uid:     comment.UId,
			Gid:     comment.GId,
			Oid:     comment.Oid,
			Content: comment.Content,
			Level:   comment.Level,
			Time:    comment.TIme.Unix(),
			Uname:   comment.UserName,
		})
	}
	response.Code = misc.CodeSuccess
	response.CommentList = commentList
	return response, nil
}

func (r *RpcLogicServer) DeleteComment(ctx context.Context, req *proto.DeleteCommentRequest) (*proto.DeleteCommentResponse, error) {
	response := &proto.DeleteCommentResponse{
		Code: misc.CodeFail,
	}
	if err := dao.DeleteComment(ctx, req.Cid); err != nil {
		misc.Logger.Error("delete comment err", zap.Error(err))
		return response, err
	}
	response.Code = misc.CodeSuccess
	return response, nil
}
