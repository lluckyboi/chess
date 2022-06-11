package chess

//PositionStruct 局面结构
type PositionStruct struct {
	SdPlayer    int      `json:"SdPlayer"`    // 轮到谁走，0=红方，1=黑方
	UcpcSquares [256]int `json:"UcpcSquares"` // 棋盘上的棋子
	RoomId      string   `json:"RoomId"`      //棋手ID
}

//NewPositionStruct 初始化棋局
func NewPositionStruct() *PositionStruct {
	p := &PositionStruct{}
	if p == nil {
		return nil
	}
	return p
}

//startup 初始化棋盘
func (p *PositionStruct) startup() {
	p.SdPlayer = 0
	for sq := 0; sq < 256; sq++ {
		p.UcpcSquares[sq] = cucpcStartup[sq]
	}
}

//changeSide 交换走子方
func (p *PositionStruct) changeSide() {
	p.SdPlayer = 1 - p.SdPlayer
}

//addPiece 在棋盘上放一枚棋子
func (p *PositionStruct) addPiece(sq, pc int) {
	p.UcpcSquares[sq] = pc
}

//delPiece 从棋盘上拿走一枚棋子
func (p *PositionStruct) delPiece(sq int) {
	p.UcpcSquares[sq] = 0
}

//movePiece 搬一步棋的棋子
func (p *PositionStruct) movePiece(mv int) int {
	sqSrc := src(mv)
	sqDst := dst(mv)
	pcCaptured := p.UcpcSquares[sqDst]
	p.delPiece(sqDst)
	pc := p.UcpcSquares[sqSrc]
	p.delPiece(sqSrc)
	p.addPiece(sqDst, pc)
	return pcCaptured
}

//makeMove 走一步棋
func (p *PositionStruct) makeMove(mv int) bool {
	pcCaptured := p.movePiece(mv)
	if p.checked() {
		p.undoMovePiece(mv, pcCaptured)
		return false
	}
	p.changeSide()
	return true
}

//generateMoves 生成所有走法
func (p *PositionStruct) generateMoves(mvs []int) int {
	nGenMoves, pcSrc, sqDst, pcDst, nDelta := 0, 0, 0, 0, 0
	pcSelfSide := sideTag(p.SdPlayer)
	pcOppSide := oppSideTag(p.SdPlayer)

	for sqSrc := 0; sqSrc < 256; sqSrc++ {
		if !inBoard(sqSrc) {
			continue
		}
		// 找到一个本方棋子，再做以下判断：
		pcSrc = p.UcpcSquares[sqSrc]
		if (pcSrc & pcSelfSide) == 0 {
			continue
		}

		// 根据棋子确定走法
		switch pcSrc - pcSelfSide {
		case PieceJiang:
			for i := 0; i < 4; i++ {
				sqDst = sqSrc + ccJiangDelta[i]
				if !inFort(sqDst) {
					continue
				}
				pcDst = p.UcpcSquares[sqDst]
				if pcDst&pcSelfSide == 0 {
					mvs[nGenMoves] = move(sqSrc, sqDst)
					nGenMoves++
				}
			}
			break
		case PieceShi:
			for i := 0; i < 4; i++ {
				sqDst = sqSrc + ccShiDelta[i]
				if !inFort(sqDst) {
					continue
				}
				pcDst = p.UcpcSquares[sqDst]
				if pcDst&pcSelfSide == 0 {
					mvs[nGenMoves] = move(sqSrc, sqDst)
					nGenMoves++
				}
			}
			break
		case PieceXiang:
			for i := 0; i < 4; i++ {
				sqDst = sqSrc + ccShiDelta[i]
				if !(inBoard(sqDst) && noRiver(sqDst, p.SdPlayer) && p.UcpcSquares[sqDst] == 0) {
					continue
				}
				sqDst += ccShiDelta[i]
				pcDst = p.UcpcSquares[sqDst]
				if pcDst&pcSelfSide == 0 {
					mvs[nGenMoves] = move(sqSrc, sqDst)
					nGenMoves++
				}
			}
			break
		case PieceMa:
			for i := 0; i < 4; i++ {
				sqDst = sqSrc + ccJiangDelta[i]
				if p.UcpcSquares[sqDst] != 0 {
					continue
				}
				for j := 0; j < 2; j++ {
					sqDst = sqSrc + ccMaDelta[i][j]
					if !inBoard(sqDst) {
						continue
					}
					pcDst = p.UcpcSquares[sqDst]
					if pcDst&pcSelfSide == 0 {
						mvs[nGenMoves] = move(sqSrc, sqDst)
						nGenMoves++
					}
				}
			}
			break
		case PieceJu:
			for i := 0; i < 4; i++ {
				nDelta = ccJiangDelta[i]
				sqDst = sqSrc + nDelta
				for inBoard(sqDst) {
					pcDst = p.UcpcSquares[sqDst]
					if pcDst == 0 {
						mvs[nGenMoves] = move(sqSrc, sqDst)
						nGenMoves++
					} else {
						if (pcDst & pcOppSide) != 0 {
							mvs[nGenMoves] = move(sqSrc, sqDst)
							nGenMoves++
						}
						break
					}
					sqDst += nDelta
				}

			}
			break
		case PiecePao:
			for i := 0; i < 4; i++ {
				nDelta = ccJiangDelta[i]
				sqDst = sqSrc + nDelta
				for inBoard(sqDst) {
					pcDst = p.UcpcSquares[sqDst]
					if pcDst == 0 {
						mvs[nGenMoves] = move(sqSrc, sqDst)
						nGenMoves++
					} else {
						break
					}
					sqDst += nDelta
				}
				sqDst += nDelta
				for inBoard(sqDst) {
					pcDst = p.UcpcSquares[sqDst]
					if pcDst != 0 {
						if (pcDst & pcOppSide) != 0 {
							mvs[nGenMoves] = move(sqSrc, sqDst)
							nGenMoves++
						}
						break
					}
					sqDst += nDelta
				}
			}
			break
		case PieceBing:
			sqDst = squareForward(sqSrc, p.SdPlayer)
			if inBoard(sqDst) {
				pcDst = p.UcpcSquares[sqDst]
				if pcDst&pcSelfSide == 0 {
					mvs[nGenMoves] = move(sqSrc, sqDst)
					nGenMoves++
				}
			}
			if hasRiver(sqSrc, p.SdPlayer) {
				for nDelta = -1; nDelta <= 1; nDelta += 2 {
					sqDst = sqSrc + nDelta
					if inBoard(sqDst) {
						pcDst = p.UcpcSquares[sqDst]
						if pcDst&pcSelfSide == 0 {
							mvs[nGenMoves] = move(sqSrc, sqDst)
							nGenMoves++
						}
					}
				}
			}
			break
		}
	}
	return nGenMoves
}

//checked 判断是否被将军
func (p *PositionStruct) checked() bool {
	nDelta, sqDst, pcDst := 0, 0, 0
	pcSelfSide := sideTag(p.SdPlayer)
	pcOppSide := oppSideTag(p.SdPlayer)

	for sqSrc := 0; sqSrc < 256; sqSrc++ {
		//找到棋盘上的帅(将)，再做以下判断：
		if !inBoard(sqSrc) || p.UcpcSquares[sqSrc] != pcSelfSide+PieceJiang {
			continue
		}

		//判断是否被对方的兵(卒)将军
		if p.UcpcSquares[squareForward(sqSrc, p.SdPlayer)] == pcOppSide+PieceBing {
			return true
		}
		for nDelta = -1; nDelta <= 1; nDelta += 2 {
			if p.UcpcSquares[sqSrc+nDelta] == pcOppSide+PieceBing {
				return true
			}
		}

		//判断是否被对方的马将军(以仕(士)的步长当作马腿)
		for i := 0; i < 4; i++ {
			if p.UcpcSquares[sqSrc+ccShiDelta[i]] != 0 {
				continue
			}
			for j := 0; j < 2; j++ {
				pcDst = p.UcpcSquares[sqSrc+ccMaCheckDelta[i][j]]
				if pcDst == pcOppSide+PieceMa {
					return true
				}
			}
		}

		//判断是否被对方的车或炮将军(包括将帅对脸)
		for i := 0; i < 4; i++ {
			nDelta = ccJiangDelta[i]
			sqDst = sqSrc + nDelta
			for inBoard(sqDst) {
				pcDst = p.UcpcSquares[sqDst]
				if pcDst != 0 {
					if pcDst == pcOppSide+PieceJu || pcDst == pcOppSide+PieceJiang {
						return true
					}
					break
				}
				sqDst += nDelta
			}
			sqDst += nDelta
			for inBoard(sqDst) {
				pcDst = p.UcpcSquares[sqDst]
				if pcDst != 0 {
					if pcDst == pcOppSide+PiecePao {
						return true
					}
					break
				}
				sqDst += nDelta
			}
		}
		return false
	}
	return false
}

//undoMovePiece 撤消搬一步棋的棋子
func (p *PositionStruct) undoMovePiece(mv, pcCaptured int) {
	sqSrc := src(mv)
	sqDst := dst(mv)
	pc := p.UcpcSquares[sqDst]
	p.delPiece(sqDst)
	p.addPiece(sqSrc, pc)
	if pcCaptured != 0 {
		p.addPiece(sqDst, pcCaptured)
	}
}

//isMate 判断是否被将死
func (p *PositionStruct) isMate() bool {
	pcCaptured := 0
	mvs := make([]int, MaxGenMoves)
	nGenMoveNum := p.generateMoves(mvs)
	for i := 0; i < nGenMoveNum; i++ {
		pcCaptured = p.movePiece(mvs[i])
		if !p.checked() {
			p.undoMovePiece(mvs[i], pcCaptured)
			return false
		}

		p.undoMovePiece(mvs[i], pcCaptured)
	}
	return true
}
