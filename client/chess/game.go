package chess

import (
	"MyChess/client/tool"
	"bytes"
	"fmt"
	"github.com/golang/freetype/truetype"
	"github.com/gorilla/websocket"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/audio"
	"github.com/hajimehoshi/ebiten/audio/wav"
	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/hajimehoshi/ebiten/text"
	"golang.org/x/image/font"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"sync"
)

//Game 象棋窗口
type Game struct {
	sqSelected     int                   // 选中的格子
	mvLast         int                   // 上一步棋
	bFlipped       bool                  //是否翻转棋盘
	bGameOver      bool                  //是否游戏结束
	showValue      string                //显示内容
	images         map[int]*ebiten.Image //图片资源
	audios         map[int]*audio.Player //音效
	audioContext   *audio.Context        //音效器
	singlePosition *PositionStruct       //棋局单例
	side           int                   //玩家是哪一方
}

var Conn *websocket.Conn
var RoomId string
var wg sync.WaitGroup
var Gmsg tool.GetTokenResp

//NewGame 创建象棋程序
func NewGame(sd int) bool {
	game := &Game{
		images:         make(map[int]*ebiten.Image),
		audios:         make(map[int]*audio.Player),
		singlePosition: NewPositionStruct(),
		side:           sd,
	}
	if game == nil || game.singlePosition == nil {
		return false
	}

	var err error
	//音效器
	game.audioContext, err = audio.NewContext(48000)
	if err != nil {
		fmt.Print(err)
		return false
	}

	//加载资源
	if ok := game.loadResource(); !ok {
		return false
	}

	//初始化棋盘 红方先走
	game.singlePosition.startup()

	//写入roomId
	game.singlePosition.RoomId = RoomId

	//更新棋盘
	if sd != -2 {
		wg.Add(1)
		go UpdateBoard(game)
	}

	//设置窗口大小
	ebiten.SetWindowSize(BoardWidth, BoardHeight)
	//标题
	ebiten.SetWindowTitle("双人象棋")
	//刷新帧数20
	ebiten.SetMaxTPS(20)

	if sd != -2 {
		if err := ebiten.RunGame(game); err != nil {
			log.Fatal(err)
			return false
		}
		wg.Wait()
	}
	return true
}

//Update 更新状态，1秒20帧 可以加载地图
func (g *Game) Update(screen *ebiten.Image) error {
	// 鼠标点一下就更新
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		x = Left + (x-BoardEdge)/SquareSize
		y = Top + (y-BoardEdge)/SquareSize
		// 点击格子时进行处理
		g.clickSquare(screen, squareXY(x, y))
	}

	g.drawBoard(screen)
	if g.bGameOver {
		g.messageBox(screen)
	}
	return nil
}

//Layout 布局采用固定尺寸即可。
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return BoardWidth, BoardHeight
}

//loadResource 加载资源
func (g *Game) loadResource() bool {
	for k, v := range resMap {
		if k >= MusicSelect {
			//加载音效
			d, err := wav.Decode(g.audioContext, audio.BytesReadSeekCloser(v))
			if err != nil {
				fmt.Print(err)
				return false
			}
			player, err := audio.NewPlayer(g.audioContext, d)
			if err != nil {
				fmt.Print(err)
				return false
			}
			g.audios[k] = player
		} else {
			//加载图片
			img, _, err := image.Decode(bytes.NewReader(v))
			if err != nil {
				fmt.Print(err)
				return false
			}
			ebitenImage, _ := ebiten.NewImageFromImage(img, ebiten.FilterDefault)
			g.images[k] = ebitenImage
		}
	}

	return true
}

//playAudio 播放音效
func (g *Game) playAudio(value int) {
	if player, ok := g.audios[value]; ok {
		player.Rewind()
		player.Play()
	}
}

//drawChess 绘制棋子
func (g *Game) drawChess(x, y int, screen, img *ebiten.Image) {
	if img == nil {
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(img, op)
}

//drawBoard 绘制棋盘
func (g *Game) drawBoard(screen *ebiten.Image) {
	//棋盘
	if v, ok := g.images[ImgChessBoard]; ok {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(0, 0)
		screen.DrawImage(v, op)
	}

	//棋子
	for x := Left; x <= Right; x++ {
		for y := Top; y <= Bottom; y++ {
			xPos, yPos := 0, 0
			if g.bFlipped {
				xPos = BoardEdge + (xFlip(x)-Left)*SquareSize
				yPos = BoardEdge + (yFlip(y)-Top)*SquareSize
			} else {
				xPos = BoardEdge + (x-Left)*SquareSize
				yPos = BoardEdge + (y-Top)*SquareSize
			}
			sq := squareXY(x, y)
			pc := g.singlePosition.UcpcSquares[sq]
			if pc != 0 {
				g.drawChess(xPos, yPos+5, screen, g.images[pc])
			}
			if sq == g.sqSelected || sq == src(g.mvLast) || sq == dst(g.mvLast) {
				g.drawChess(xPos, yPos, screen, g.images[ImgSelect])
			}
		}
	}
}

//clickSquare 点击格子处理
func (g *Game) clickSquare(screen *ebiten.Image, sq int) {
	//得到点击的格子上的棋子
	pc := 0
	//判断是否翻转
	if g.bFlipped {
		pc = g.singlePosition.UcpcSquares[squareFlip(sq)]
	} else {
		pc = g.singlePosition.UcpcSquares[sq]
	}
	//检查是否轮到自己
	if g.singlePosition.SdPlayer == g.side {
		//sideTag 红方为8 黑方为16
		if pc&sideTag(g.singlePosition.SdPlayer) != 0 {
			//如果点击自己的棋子，那么直接选中
			g.sqSelected = sq
			g.playAudio(MusicSelect)
		} else if g.sqSelected != 0 && !g.bGameOver {
			//如果点击的不是自己的棋子，但有棋子选中了(一定是自己的棋子,比如吃子)，那么走这个棋子
			mv := move(g.sqSelected, sq)
			if g.singlePosition.legalMove(mv) {
				//如果没有将军
				if g.singlePosition.makeMove(mv) {
					g.mvLast = mv
					g.sqSelected = 0
					if g.singlePosition.isMate() {
						// 如果分出胜负，那么播放胜负的声音，并且弹出不带声音的提示框
						g.playAudio(MusicGameWin)
						tool.AddWinCount(Gmsg.Token)
						g.bGameOver = true
						//ws 更新棋盘
						UpdateOtherBoard(g)
					} else {
						// 如果没有分出胜负，那么播放将军、吃子或一般走子的声音
						if g.singlePosition.checked() {
							g.playAudio(MusicJiang)
						} else {
							if pc != 0 {
								g.playAudio(MusicEat)
							} else {
								g.playAudio(MusicPut)
							}
						}
						//ws 更新棋盘
						UpdateOtherBoard(g)
					}
				} else {
					g.playAudio(MusicJiang) // 播放被将军的声音
				}
			}
			//如果根本就不符合走法(例如马不走日字)，那么不做任何处理
		}
	}
	//如果没轮到自己，什么都不干
}

//messageBox 提示
func (g *Game) messageBox(screen *ebiten.Image) {
	if g.side == (1 - g.singlePosition.SdPlayer) {
		g.showValue = "You Win!"
	} else {
		g.showValue = "You Lose!"
	}
	fmt.Println(g.showValue)
	tt, err := truetype.Parse(fonts.ArcadeN_ttf)
	if err != nil {
		fmt.Print(err)
		return
	}
	arcadeFont := truetype.NewFace(tt, &truetype.Options{
		Size:    16,
		DPI:     72,
		Hinting: font.HintingFull,
	})

	text.Draw(screen, g.showValue, arcadeFont, 180, 288, color.Black)
}
