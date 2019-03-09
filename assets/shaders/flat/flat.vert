#version 410 core

layout(location = 0) in vec3 pos;

uniform mat4 M, V, P;
uniform vec3 flatColor = vec3(1, 1, 1);

out vec3 col0;

void main(){
    gl_Position = P * V * M * vec4(pos, 1.0);
    col0 = flatColor;
}