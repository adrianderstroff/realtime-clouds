#version 430

layout (location = 0) in vec2 position;

out VertexIn {
    vec2 pos;
    int id;
    int vid;
} o;

void main() {
    gl_Position = vec4(0.0, 0.0, 0.0, 1.0);
    o.pos = position;
    o.id  = gl_InstanceID;
    o.vid = gl_VertexID;
}