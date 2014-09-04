package main

import (
  "bufio"
  "code.google.com/p/draw2d/draw2d"
  "fmt"
  "image"
  "image/draw"
  "image/color"
  "image/png"
  "os"
  "math"
  "github.com/adammck/ik"
)

type Projection struct {
  img          draw.Image
  canvasWidth  float64
  canvasHeight float64
  worldWidth   float64
  worldHeight  float64
}

var (
  grey   = color.RGBA{0xCC, 0xCC, 0xCC, 0xFF}
  black  = color.RGBA{0,    0,    0,    0xFF}
  red    = color.RGBA{0xAA, 0,    0,    0xFF}
  green  = color.RGBA{0,    0xAA, 0,    0xFF}
  blue   = color.RGBA{0,    0,    0xAA, 0xFF}
  yellow = color.RGBA{0xAA, 0xAA, 0,    0xFF}
  cyan   = color.RGBA{0,    0xAA, 0xAA, 0xFF}
  purple = color.RGBA{0xAA, 0,    0xAA, 0xFF}
)

func MakeProjection(cw int, ch int, ww float64, wh float64) *Projection {
  img := image.NewRGBA(image.Rect(0, 0, cw, ch))
  draw.Draw(img, img.Bounds(), image.White, image.ZP, draw.Src)
  return &Projection{img, float64(cw), float64(ch), ww, wh}
}

func main() {

  target := ik.MakeVector3(25, -10, 0)
  step := (math.Pi/180) * 4.5

  x := ik.MakeRootSegment(ik.MakeVector3(5, 0, 0))
  a := ik.MakeSegment(x, ik.Euler(0, 0,  -18), ik.Euler(0, 0, 72), ik.MakeVector3(20, 0, 0))
  b := ik.MakeSegment(a, ik.Euler(0, 0, -135), ik.Euler(0, 0, 45), ik.MakeVector3(10, 0, 0))
  c := ik.MakeSegment(b, ik.Euler(0, 0,  -90), ik.Euler(0, 0, 45), ik.MakeVector3(5, 0, 0))

  bestDistance := math.Inf(1)
  bestAngles := [3]ik.EulerAngles{}
  n := 0

  p := MakeProjection(1000, 1000, 100.0, 100.0)
  p.grid(5, 5)

  // TODO: Move all this shit into Segment or some kind of Solver
  for _, aea := range a.Range(step) {
    for _, bea := range b.Range(step) {
      for _, cea := range c.Range(step) {

        a.SetRotation(aea)
        b.SetRotation(bea)
        c.SetRotation(cea)
        n++

        distanceFromTarget := target.Distance(c.End())
        if bestDistance > distanceFromTarget {
          bestDistance = distanceFromTarget
          bestAngles = [3]ik.EulerAngles{
            *aea,
            *bea,
            *cea,
          }
        }

        cp := c.End()
        p.cross(cp.X, cp.Y, grey)
      }
    }
  }

  fmt.Printf("checked %d positions\n",n)
  fmt.Printf("distance from target: %0.4f\n",bestDistance)
  fmt.Printf("segment angles: %v\n", bestAngles)

  // Restore the best Rotation
  a.SetRotation(&bestAngles[0])
  b.SetRotation(&bestAngles[1])
  c.SetRotation(&bestAngles[2])
  p.drawSegment(x, green)

  p.cross(target.X, target.Y, red)
  write(p.img, "image.png")
}

// Projects an (x,y) world coordinate onto the canvas, returning (x,y).
func (p *Projection) project(wx float64, wy float64) (float64, float64) {
  return ((wx / p.worldWidth) * p.canvasWidth) + (p.canvasWidth * 0.5),
         ((wy / p.worldHeight) * -p.canvasHeight) + (p.canvasHeight * 0.5)
}

// drawSegment draws the specified limb on the canvas, by starting at the given
// object, and recursing for each child.
func (p *Projection) drawSegment(s *ik.Segment, col color.RGBA) {
  a := s.Start()
  b := s.End()
  p.line(a.X, a.Y, b.X, b.Y, col)

  if s.Child != nil {
    p.drawSegment(s.Child, col)
  }
}

// Cross draws a cross on the canvas, at the specified world coordinates.
func (p *Projection) cross(wx float64, wy float64, col color.RGBA) {
  cx, cy := p.project(wx, wy)

  c := draw2d.NewGraphicContext(p.img)
  c.SetStrokeColor(col)
  c.SetLineWidth(1)
  size := 2.0

  // top left -> bottom right
  c.MoveTo(cx - size, cy - size)
  c.LineTo(cx + size, cy + size)
  c.Stroke()

  // top right -> bottom left
  c.MoveTo(cx + size, cy - size)
  c.LineTo(cx - size, cy + size)
  c.Stroke()
}

// Line draws a line on the canvas, between the specified world coordinates.
func (p *Projection) line(w1x float64, w1y float64, w2x float64, w2y float64, col color.RGBA) {
  c1x, c1y := p.project(w1x, w1y)
  c2x, c2y := p.project(w2x, w2y)

  c := draw2d.NewGraphicContext(p.img)
  c.SetStrokeColor(col)
  c.SetLineWidth(1)

  c.MoveTo(c1x, c1y)
  c.LineTo(c2x, c2y)
  c.Stroke()
}

// Grid renders a grid onto the canvas.
func (p *Projection) grid(xInterval float64, yInterval float64) {
  l := draw2d.NewGraphicContext(p.img)
  l.SetStrokeColor(color.RGBA{0xEE, 0xEE, 0xEE, 0xFF})
  l.SetLineWidth(0.5)

  xCount := p.worldWidth / xInterval
  yCount := p.worldHeight / yInterval

  // horizontal lines
  for x := 1.0; x < xCount; x += 1 {
    xx, _ := p.project((x - (xCount / 2)) * xInterval, 0)
    l.MoveTo(xx, 0)
    l.LineTo(xx, p.canvasHeight)
    l.Stroke()
  }

  // vertical lines
  for y := 1.0; y < yCount; y += 1 {
    _, yy := p.project(0, (y - (yCount / 2)) * yInterval)
    l.MoveTo(0, yy)
    l.LineTo(p.canvasWidth, yy)
    l.Stroke()
  }

    l.SetStrokeColor(color.RGBA{0xAA, 0xAA, 0xAA, 0xFF})

    // horiz axis
  l.MoveTo(p.canvasWidth/2, 0)
  l.LineTo(p.canvasWidth/2, p.canvasHeight)
  l.Stroke()

  // vert axis
  l.MoveTo(0, p.canvasHeight/2)
  l.LineTo(p.canvasWidth, p.canvasHeight/2)
  l.Stroke()
}

func write(img draw.Image, path string) {
  f, err := os.Create(path)
  defer f.Close()
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }

  buf := bufio.NewWriter(f)
  err = png.Encode(buf, img)
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }

  err = buf.Flush()
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }

  fmt.Printf("Wrote: %s\n", path)
}
