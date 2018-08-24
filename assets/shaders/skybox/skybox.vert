#version 400
layout(location = 0) in vec3 pos;

uniform mat4 M, V, P;

out vec3 tex;

void main(){
    gl_Position = P*V*M*vec4(pos, 1.0);
    //gl_Position = (P*V*M*vec4(pos, 1.0)).xyww;
    //gl_Position.z -= 0.00001;
    tex = normalize(pos);
}