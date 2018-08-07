#version 430

in  VertexIn {
    int id;
} i[];
out VertexOut {
    vec3 position;
    vec3 normal;
} o;
struct Tile {
    vec4  tri1;
    vec4  tri2;
    vec2  pos;
    float lod;
    float padding;
};

layout(points) in;
layout(triangle_strip, max_vertices = 6) out;
layout(std430, binding = 0) buffer TileBuffer { Tile tiles[]; };

uniform mat4 M, V, P;
uniform float tilesize;

vec3 calcHeight(vec4 plane, float x, float z) {
    float y = -(plane.w + plane.x*x + plane.z*z) / plane.y;
    return vec3(0, y, 0);
}

void main() {
    // setup vectors
    vec3 center = gl_in[0].gl_Position.xyz;
    vec3 dx     = vec3(tilesize/2, 0, 0);
    vec3 dz     = vec3(0, 0, tilesize/2);

    // get all positions
    vec3 p1     = center - dx + dz;
    vec3 p2     = center - dx - dz;
    vec3 p3     = center + dx + dz;
    vec3 p4     = center + dx - dz;

    // get height
    Tile tile = tiles[i[0].id];
    p1 = p1 + calcHeight(tile.tri1, p1.x, p1.z);
    p2 = p2 + calcHeight(tile.tri1, p2.x, p2.z);
    p3 = p3 + calcHeight(tile.tri1, p3.x, p3.z);
    p4 = p4 + calcHeight(tile.tri2, p4.x, p4.z);

    // get normals
    vec3 n1 = tile.tri1.xyz;
    vec3 n2 = tile.tri2.xyz;

    // create triangle 1
    o.position = p1;
    o.normal   = n1;
    gl_Position = P*V*M * vec4(p1, 1.0);
    EmitVertex();
    o.position = p2;
    o.normal   = n1;
    gl_Position = P*V*M * vec4(p2, 1.0);
    EmitVertex();
    o.position = p3;
    o.normal   = n1;
    gl_Position = P*V*M * vec4(p3, 1.0);
    EmitVertex();
    EndPrimitive();

    // create triangle 2
    o.position = p3;
    o.normal   = n2;
    gl_Position = P*V*M * vec4(p3, 1.0);
    EmitVertex();
    o.position = p2;
    o.normal   = n2;
    gl_Position = P*V*M * vec4(p2, 1.0);
    EmitVertex();
    o.position = p4;
    o.normal   = n2;
    gl_Position = P*V*M * vec4(p4, 1.0);
    EmitVertex();
    EndPrimitive();  
}