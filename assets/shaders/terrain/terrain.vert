#version 430

layout (location = 0) in vec3 position;

out VertexIn {
    int id;
} o;

void main() {
    gl_Position = vec4(position, 1.0);
    o.id = gl_VertexID;
}