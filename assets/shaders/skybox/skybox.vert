#version 400
layout(location = 0) in vec3 vert;

uniform mat4 M, V, P;

out vec3 tex;

void main(){
    gl_Position = (P*V*vec4(vert, 1.0)).xyww;
    gl_Position.z -= 0.000001;
    tex = normalize(vert);
}