# realtime-clouds

Implementation of the volumetric cloud rendering approach based on the presentation ***The real-time volumetric cloudscapes of Horizon:Zero Dawn*** by Andrew Schneider held at ***SIGGRAPH 2015*** as part of the *Advances in real-time rendering* course.

## Requirements
This project requires a GPU with OpenGL 4.3+ support.

The following dependencies depend on cgo. To make them work under Windows a compatible version of **mingw** is necessary. Information can be found [here](https://github.com/go-gl/glfw/issues/91). In my case I used *x86_64-7.2.0-posix-seh-rt_v5-rev1*. After installing the right version of **mingw** you can continue by installing the dependencies that follow next.

This project depends on **glfw** for creating a window and providing a rendering context, **go-gl/gl** for providing bindings to OpenGL and **go-gl/mathgl** provides vector and matrix math for OpenGL.
```
go get -u github.com/go-gl/glfw/v3.2/glfw
go get -u github.com/go-gl/gl/v4.3-core/gl
go get -u github.com/go-gl/mathgl/mgl32
```
After getting all dependencies the project should work without any errors.

## Style guide

### Dimensions 
*length x width x height*

## Theory

TODO