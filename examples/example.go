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
  "github.com/adammck/ik"
)

type Axis int

type Projection struct {
  label        string
  horizAxis    Axis
  vertiAxis    Axis
  offsetLeft   float64
  offsetTop    float64
  canvasWidth  float64
  canvasHeight float64
  worldWidth   float64
  worldHeight  float64
}

type ProjectionSet struct {
  img  draw.Image
  proj []*Projection
}

const (
  X Axis = iota
  Y Axis = iota
  Z Axis = iota
)


var (
  grey      = color.RGBA{0xCC, 0xCC, 0xCC, 0xFF}
  lightGrey = color.RGBA{0xEE, 0xEE, 0xEE, 0xFF}
  black     = color.RGBA{0,    0,    0,    0xFF}
  red       = color.RGBA{0xAA, 0,    0,    0xFF}
  green     = color.RGBA{0,    0xAA, 0,    0xFF}
  blue      = color.RGBA{0,    0,    0xAA, 0xFF}
  yellow    = color.RGBA{0xAA, 0xAA, 0,    0xFF}
  cyan      = color.RGBA{0,    0xAA, 0xAA, 0xFF}
  purple    = color.RGBA{0xAA, 0,    0xAA, 0xFF}
)

func main() {

  target := ik.MakeVector3(30, -15, 10)
  fmt.Printf("target: %v\n", target)

  // The position of the object in space must be specified by two segments. The
  // first positions it, then the second (which is always zero-length) rotates
  // it into the home orientation.
  r1 := ik.MakeRootSegment(*ik.MakeVector3(10, 0, 10))
  r2 := ik.MakeSegment(r1, ik.Euler(0, 0, 0), ik.Euler(0, 0, 0), *ik.MakeVector3(0, 0, 0)) // -60 0 0

  // Movable segments
  coxa   := ik.MakeSegment(r2,    ik.Euler(60, 0,   0), ik.Euler(-60, 0,    0), *ik.MakeVector3( 5, -5, 0)) // 0 0 0
  femur  := ik.MakeSegment(coxa,  ik.Euler(0,  0,  90), ik.Euler(  0, 0,    0), *ik.MakeVector3(10,  0, 0)) // 0 0 60
  tibia  := ik.MakeSegment(femur, ik.Euler(0,  0, -45), ik.Euler(  0, 0, -135), *ik.MakeVector3(10,  0, 0)) // 0 0 -90
  tarsus := ik.MakeSegment(tibia, ik.Euler(0,  0,   0), ik.Euler(  0, 0,  -90), *ik.MakeVector3( 5,  0, 0)) // 0 0 -60
  _ = tarsus

  img := image.NewRGBA(image.Rect(0, 0, 1000, 1000))
  draw.Draw(img, img.Bounds(), image.White, image.ZP, draw.Src)

  top := &Projection{
    label:        "Top (x,z)",
    horizAxis:    X,
    vertiAxis:    Z,
    offsetLeft:   0,
    offsetTop:    0,
    canvasWidth:  500,
    canvasHeight: 500,
    worldWidth:   100.0,
    worldHeight:  100.0,
  }

  front := &Projection{
    label:        "Front (x,y)",
    horizAxis:    X,
    vertiAxis:    Y,
    offsetLeft:   0,
    offsetTop:    500,
    canvasWidth:  500,
    canvasHeight: 500,
    worldWidth:   100.0,
    worldHeight:  100.0,
  }

  side := &Projection{
    label:        "Side (z,y)",
    horizAxis:    Z,
    vertiAxis:    Y,
    offsetLeft:   500,
    offsetTop:    500,
    canvasWidth:  500,
    canvasHeight: 500,
    worldWidth:   100.0,
    worldHeight:  100.0,
  }

  p := &ProjectionSet{img, []*Projection{top, front, side}}
  p.drawGrids(lightGrey, 5)
  p.drawLabels(black)
  p.drawSplits(black)

  best := ik.Solve(coxa, target, 0.1, func(v ik.Vector3, d float64) {
    p.cross(v, grey)
  })

  fmt.Printf("\n\n\n\ndistance: %0.4f\n", best.Distance)
  fmt.Printf("segment: %s\n", best.Segment)
  p.drawSegment(best.Segment, red)

  p.cross(ik.Vector3{0, 0, 0}, blue)
  p.cross(*target, blue)

  write(p.img, "image.png")
}

// Projects an (x,y) world coordinate onto the canvas, returning (x,y).
func (p *Projection) project(v ik.Vector3) (float64, float64) {
  var x, y float64

  switch p.horizAxis {
  case X:
    x = v.X
  case Y:
    x = v.Y
  case Z:
    x = v.Z
  default:
    panic("invalid horizAxis")
  }

  switch p.vertiAxis {
  case X:
    y = v.X
  case Y:
    y = v.Y
  case Z:
    y = v.Z
  default:
    panic("invalid vertiAxis")
  }

  return p.offsetLeft + ((x / p.worldWidth) * p.canvasWidth) + (p.canvasWidth * 0.5),
         p.offsetTop + ((y / p.worldHeight) * -p.canvasHeight) + (p.canvasHeight * 0.5)
}


func (ps *ProjectionSet) cross(v ik.Vector3, col color.RGBA) {
  for _, p := range ps.proj {
    p.cross(ps.img, v, col)
  }
}

// drawSegment draws the specified limb on the canvas, by starting at the given
// object, and recursing for each child.
func (ps *ProjectionSet) drawSegment(s *ik.Segment, col color.RGBA) {
  for _, p := range ps.proj {
    p.drawSegment(ps.img, s, col)
  }
}


func (ps *ProjectionSet) drawLabels(col color.RGBA) {
  c := draw2d.NewGraphicContext(ps.img)

  for _, p := range ps.proj {
    draw2d.SetFontFolder("/Users/adammck/code/go/src/code.google.com/p/draw2d/resource/font")
    c.SetFontData(draw2d.FontData{"luxi", draw2d.FontFamilySans, draw2d.FontStyleNormal})
    c.SetFillColor(col)
    c.SetFontSize(10)
    c.FillStringAt(p.label, p.offsetLeft + 5, p.offsetTop + 15)
  }
}

func (ps *ProjectionSet) drawSplits(col color.RGBA) {
  c := draw2d.NewGraphicContext(ps.img)
  c.SetStrokeColor(col)
  c.SetLineWidth(2)

  // Horiz
  c.MoveTo(0, 500)
  c.LineTo(1000, 500)
  c.Stroke()

  // Verti
  c.MoveTo(500, 0)
  c.LineTo(500, 1000)
  c.Stroke()
}

func (ps *ProjectionSet) drawGrids(col color.RGBA, interval float64) {
  for _, p := range ps.proj {
    p.grid(ps.img, interval, col)
  }
}


func (p *Projection) drawSegment(img draw.Image, s *ik.Segment, col color.RGBA) {
  p.line(img, s.Start(), s.End(), col)
  p.cross(img, s.End(), black)

  if s.Child != nil {
    p.drawSegment(img, s.Child, col)
  }
}

// Cross draws a cross on the canvas, at the specified world coordinates.
func (p *Projection) cross(img draw.Image, v ik.Vector3, col color.RGBA) {
  cx, cy := p.project(v)

  c := draw2d.NewGraphicContext(img)
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
func (p *Projection) line(img draw.Image, v1 ik.Vector3, v2 ik.Vector3, col color.RGBA) {
  c1x, c1y := p.project(v1)
  c2x, c2y := p.project(v2)

  c := draw2d.NewGraphicContext(img)
  c.SetStrokeColor(col)
  c.SetLineWidth(1)

  c.MoveTo(c1x, c1y)
  c.LineTo(c2x, c2y)
  c.Stroke()
}

// Grid renders a grid onto the specified projection.
func (p *Projection) grid(img draw.Image, interval float64, col color.RGBA) {
  for x := (-p.worldWidth * 0.5) + interval; x < (p.worldWidth * 0.5); x += interval {
    for y := (-p.worldHeight * 0.5) + interval; y < (p.worldHeight * 0.5); y += interval {
      v := ik.Vector3{0, 0, 0}

      switch p.horizAxis {
      case X:
        v.X = x
      case Y:
        v.Y = x
      case Z:
        v.Z = x
      default:
        panic("invalid horizAxis")
      }

      switch p.vertiAxis {
      case X:
        v.X = y
      case Y:
        v.Y = y
      case Z:
        v.Z = y
      default:
        panic("invalid vertiAxis")
      }

      p.cross(img, v, col)
    }
  }
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
