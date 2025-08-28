# Render2Go

A new generation engine aimed at rendering animation.

USING GOLANG!

## Features

- **Dndependent less**: you can run it out of golang env.
- **Use a simply language**: low learning price. 
- **Color themes**: support color management.

## Usage

### Generate binary file in your platform

```bash
# windows:
go build -o render2go.exe ./cmd/render2go
#ubuntu:
go build -o render2go ./cmd/render2go
```

*you can download pre-compiled file on the download page*

### Use it!

```bash
render2go.exe script.r2g

render2go.exe -i

render2go.exe -debug script.r2g

render2go.exe --help
```

## Basic Grammar

### SCENE
```r2g
scene 800 600 "my_project"
```

### CREATE
```r2g
TEXT
create text title "Hello World" 24 (400, 100)

CIRCLE
create circle shape 50 (400, 300)

RECTANGLE
create rectangle box 100 80 (200, 200)

LINE
create line connector (100, 100) (200, 200)

ARROW
create arrow pointer (150, 150) (250, 250)

POLYGON
create polygon star [(300, 160), (320, 200), (360, 200)]
```

### SET
```r2g
set title.color = lightpurple
set shape.color = midblue
set shape.opacity = 0.8
```

### RENDER
```r2g
render
save "output.png"
```

## COLOR MANAGEMENT

- `darkcolor` - #0a2639
- `purpleblue` - #196090
- `midblue` - #3498db
- `cyanblue` - #8bc4ea
- `lightpurple` - #d4e9f7

## STRACTURE

```
render2go/
├── cmd/render2go/           
├── interpreter/           
│   ├── lexer.go          
│   ├── parser.go         
│   ├── evaluator.go      
│   └── interpreter.go    
├── scripts/              
│   ├── main.r2g         
│   ├── simple.r2g      
│   ├── advanced.r2g     
│   └── showcase.r2g     
├── output/               
└── examples/            
```

## EXAMPLES

### 1
```gma
scene 400 300 "hello_world"

create text greeting "Hello Render2Go!" 24 (200, 150)
set greeting.color = lightpurple

render
save "hello.png"
```

### 2
```gma
scene 600 400 "complex_demo"

create text title "Advanced Demo" 28 (300, 40)
set title.color = lightpurple

create circle c1 25 (150, 200)
set c1.color = midblue
set c1.opacity = 0.9

create rectangle r1 50 35 (300, 200)
set r1.color = cyanblue
set r1.opacity = 0.8

create arrow a1 (200, 200) (250, 200)
set a1.color = lightpurple

render
save "complex.png"
```

---

**Render2Go**
