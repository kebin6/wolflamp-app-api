package game

import (
	"context"
	"github.com/kebin6/wolflamp-app-api/internal/logic/player"
	"github.com/kebin6/wolflamp-rpc/types/wolflamp"
	"github.com/suyuan32/simple-admin-common/i18n"
	"github.com/zeromicro/go-zero/core/errorx"
	"google.golang.org/grpc/status"

	"github.com/kebin6/wolflamp-app-api/internal/svc"
	"github.com/kebin6/wolflamp-app-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRoundInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetRoundInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRoundInfoLogic {
	return &GetRoundInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx}
}

func (l *GetRoundInfoLogic) GetRoundInfo(req *types.RoundReq) (resp *types.RoundResp, err error) {

	if !l.svcCtx.Config.WolfLampRpc.Enabled {
		return nil, errorx.NewCodeUnavailableError(i18n.ServiceUnavailable)
	}

	id, err := player.NewPersonInfoLogic(l.ctx, l.svcCtx).GetPlayerId()
	if err != nil {
		return nil, err
	}

	// 获取当前回合数据
	round, err := l.svcCtx.WolfLampRpc.FindRound(l.ctx, &wolflamp.FindRoundReq{Id: req.Id})
	if err != nil {
		if status.Convert(err).Message() == "target does not exist" {
			return nil, errorx.NewCodeInvalidArgumentError("game.roundNotFound")
		}
		return nil, err
	}

	// 获取当前回合所有投注数据
	invests, err := l.svcCtx.WolfLampRpc.GetInvestByRoundId(l.ctx, &wolflamp.GetInvestsByRoundIdReq{RoundId: round.Id})
	if err != nil {
		return nil, err
	}
	foldInvestPlayerNum := make(map[uint32]uint32)
	// 以羊圈纬度统计各羊圈投注玩家数
	for _, invest := range invests.Data {
		if _, ok := foldInvestPlayerNum[invest.LambNum]; ok {
			foldInvestPlayerNum[invest.FoldNo]++
		} else {
			foldInvestPlayerNum[invest.FoldNo] = 1
		}
	}

	folds := types.FoldInfo{}
	if round.Folds != nil {
		for _, fold := range round.Folds {
			playerNum := uint32(0)
			if _, ok := foldInvestPlayerNum[fold.FoldNo]; ok {
				playerNum = foldInvestPlayerNum[fold.FoldNo]
			}
			detail := types.FoldDetail{
				FoldNo:    fold.FoldNo,
				LambNum:   fold.LambNum,
				PlayerNum: playerNum,
			}
			switch fold.FoldNo {
			case 1:
				folds.Fold1 = detail
			case 2:
				folds.Fold2 = detail
			case 3:
				folds.Fold3 = detail
			case 4:
				folds.Fold4 = detail
			case 5:
				folds.Fold5 = detail
			case 6:
				folds.Fold6 = detail
			case 7:
				folds.Fold7 = detail
			case 8:
				folds.Fold8 = detail
			}
		}
	}

	// 获取上一回合被选中的羊圈
	previousRound, err := l.svcCtx.WolfLampRpc.PreviousRound(l.ctx, &wolflamp.Empty{})
	if err != nil {
		return nil, err
	}

	// 获取当前玩家投注情况
	investResp, err := l.svcCtx.WolfLampRpc.GetInvestInfoByPlayerId(l.ctx,
		&wolflamp.GetInvestInfoByPlayerIdReq{PlayerId: *id, RoundId: &round.Id})
	if err != nil {
		return nil, err
	}
	playerInvestInfo := types.PlayerInvestInfo{
		Fold1: types.FoldDetail{FoldNo: 1, LambNum: 0},
		Fold2: types.FoldDetail{FoldNo: 2, LambNum: 0},
		Fold3: types.FoldDetail{FoldNo: 3, LambNum: 0},
		Fold4: types.FoldDetail{FoldNo: 4, LambNum: 0},
		Fold5: types.FoldDetail{FoldNo: 5, LambNum: 0},
		Fold6: types.FoldDetail{FoldNo: 6, LambNum: 0},
		Fold7: types.FoldDetail{FoldNo: 7, LambNum: 0},
		Fold8: types.FoldDetail{FoldNo: 8, LambNum: 0},
	}
	if investResp.Data != nil {
		for _, invest := range investResp.Data {
			switch invest.FoldNo {
			case 1:
				playerInvestInfo.Fold1.LambNum += invest.LambNum
			case 2:
				playerInvestInfo.Fold2.LambNum += invest.LambNum
			case 3:
				playerInvestInfo.Fold3.LambNum += invest.LambNum
			case 4:
				playerInvestInfo.Fold4.LambNum += invest.LambNum
			case 5:
				playerInvestInfo.Fold5.LambNum += invest.LambNum
			case 6:
				playerInvestInfo.Fold6.LambNum += invest.LambNum
			case 7:
				playerInvestInfo.Fold7.LambNum += invest.LambNum
			case 8:
				playerInvestInfo.Fold8.LambNum += invest.LambNum
			}
		}
	}

	// 获取开奖结果
	resultInfo := types.ResultInfo{
		FoldNum:       round.SelectedFold,
		ProfitAndLoss: 0,
	}
	if round.SelectedFold != 0 {
		playerInvests, err := l.svcCtx.WolfLampRpc.GetInvestInfoByPlayerId(l.ctx,
			&wolflamp.GetInvestInfoByPlayerIdReq{RoundId: &round.Id, PlayerId: *id})
		if err != nil {
			return nil, err
		}
		profitAndLoss := float32(0)
		for _, invest := range playerInvests.Data {
			profitAndLoss += invest.ProfitAndLoss
		}
		resultInfo.ProfitAndLoss = profitAndLoss
	}

	return &types.RoundResp{
		Data: types.RoundInfo{
			Id:               round.Id,
			FoldInfo:         folds,
			PreviousFoldNo:   previousRound.SelectedFoldNo,
			PlayerInvestInfo: playerInvestInfo,
			Status:           round.Status,
			StartAt:          round.StartAt,
			OpenAt:           round.OpenAt,
			EndAt:            round.EndAt,
			ResultInfo:       resultInfo,
		},
	}, nil

}
