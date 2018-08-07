#version 410 core

layout(location = 0) in vec3 vert;

void main(){
    gl_Position = vec4(vert.xy, 0.0, 1.0);
}