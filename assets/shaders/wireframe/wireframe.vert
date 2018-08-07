#version 410 core

layout(location = 0) in vec3 position;
layout(location = 1) in vec3 barycoord;

uniform mat4 M, V, P;

out Vertex {
    vec3 barycoord;
} o;

void main(){
    o.barycoord = barycoord;
    gl_Position = P * V * M * vec4(position, 1.0);
}