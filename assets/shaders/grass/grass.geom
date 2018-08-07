#version 430

//-----------------------------------------------------------------------------------//
// constants                                                                         //
//-----------------------------------------------------------------------------------//
const float TWOPI  = 6.28318530717;
const float PI     = 3.14159265358;
const float HALFPI = 1.57079632679;

//-----------------------------------------------------------------------------------//
// data structs                                                                      //
//-----------------------------------------------------------------------------------//
in  VertexIn {
    vec2 pos;
    int  id;
    int  vid;
} i[];
out VertexOut {
    vec3  position;
    vec2  uv;
    vec3  normal;
    float texID;
} o;
struct Tile {
    vec4  tri1;
    vec4  tri2;
    vec2  pos;
    float lod;
    float padding;
};

//-----------------------------------------------------------------------------------//
// in out data                                                                       //
//-----------------------------------------------------------------------------------//
layout(points) in;
layout(triangle_strip, max_vertices = 20) out;

//-----------------------------------------------------------------------------------//
// buffers                                                                           //
//-----------------------------------------------------------------------------------//
layout(std430, binding = 0) buffer TileBuffer    { Tile tiles[]; };
layout(std430, binding = 1) buffer Velocityfield { vec4 velocity[]; };

//-----------------------------------------------------------------------------------//
// uniforms                                                                          //
//-----------------------------------------------------------------------------------//
uniform mat4  M, V, P;
uniform vec3  cameraPos;
uniform float grassHeight;
uniform int   bladeCount;
uniform float tilesize;
uniform float t;
uniform float d2;
uniform int   radius;

//-----------------------------------------------------------------------------------//
// get tile                                                                          //
//-----------------------------------------------------------------------------------//
Tile getTile() {
    return tiles[i[0].id];
}
vec3 getTilePos() {
    vec2 pos = getTile().pos;
    return vec3(pos.x, 0, pos.y);
}

//-----------------------------------------------------------------------------------//
// randomization                                                                     //
//-----------------------------------------------------------------------------------//
float rand() {
    Tile tile = getTile();
    float f = tile.pos.x*tile.pos.y;
    return sin(f*HALFPI*fract(i[0].pos.x) + f*HALFPI*fract(i[0].pos.y));
}
float range(float min, float max) {
    return (max - min)*rand() + min;
}
float time(float freq, float phase) {
    return sin(phase + mod(t,freq)/(freq*0.5) * PI);
}

//-----------------------------------------------------------------------------------//
// get the root positions                                                            //
//-----------------------------------------------------------------------------------//
vec3 calcRootHeight(vec4 plane, float x, float z) {
    float y = -(plane.w + plane.x*x + plane.z*z) / plane.y;
    return vec3(0, y, 0);
}
vec3 getRootLocalPos(float r) {
    vec2 grassPos = i[0].pos;
    grassPos.x = mod(grassPos.x + r, 1.0) - 0.5;
    grassPos.y = mod(grassPos.y + r, 1.0) - 0.5;
    return vec3(grassPos.x, 0.0, grassPos.y);
}
vec3 calcRootWorldPos(vec3 local) {
    // get grass local coordinates
    vec3 rPos = local*tilesize;

    // get grass world coordinates
    vec3 pos = getTilePos() + rPos;

    // add root height
    Tile tile = getTile();
    vec4 plane = tile.tri2;
    if (rPos.x < rPos.z) { plane = tile.tri1; }
    pos = pos + calcRootHeight(plane, pos.x, pos.z);

    return pos;
}
vec3 extractCameraRight() {
    return vec3(V[0][0], V[1][0], V[2][0]);
}

//-----------------------------------------------------------------------------------//
// calculating LOD                                                                   //
//-----------------------------------------------------------------------------------//
float dist(vec3 pos) {
    return length(pos - cameraPos);
}
int   calcLOD(vec3 pos) {
    float x = dist(pos) / d2;
    float off = 0.5;
    return int(4*(off*off)/((x+off)*(x+off)));
}
float calcLODDist(vec3 pos) {
    float x = 5 * dist(pos) / d2;
    return pow(1.17, -(x*x));
}
float calcLODBladeHeight(vec3 pos) {
    return calcLODDist(pos);
}
float calcLODBladeWidth(vec3 pos) {
    float dist = (1-calcLODDist(pos))*4 + 1.0;
    return range(0.6, 0.9)*dist;
}
int   calcLODBladeCount(vec3 pos) {
    int lod = calcLOD(pos);
    int count = bladeCount;
    if     (lod == 2) count = int(0.95*bladeCount);
    else if(lod == 1) count = int(0.85*bladeCount);
    else if(lod == 0) count = 1;
    return count;
}

//-----------------------------------------------------------------------------------//
// creating the grass segments                                                       //
//-----------------------------------------------------------------------------------//
vec3 calcNormal(vec3 v1, vec3 v2, vec3 v3) {
    vec3 d1 = v1 - v2;
    vec3 d2 = v3 - v2;
    return normalize(cross(d1, d2));
}
vec3 displace(vec3 pos, vec3 root, vec3 wind, float coeff, float r) {
    // calc wind influence
    vec3 force = range(0.4, 0.5)*wind;

    // calc grass bending
    vec3 npos = pos; 
    npos.xz += force.xz*coeff;
    npos.y  -= length(force)*sqrt(coeff);

    // calc grass oszillation
    vec3 dir = normalize(force);
    float strength = length(force);
    if(strength == 0.0) { dir = vec3(1, 0, 0); }
    float bend = 1 - (grassHeight - (npos.y - root.y)) / grassHeight;
    bend = 0.4*strength * bend*bend * time(60 + r*20, r*PI);
    npos.xz += bend*dir.xz*coeff;
    npos.y  -= bend*sqrt(coeff);

    return npos;
}
void makeSegment(vec3 root, vec3 right, float s, float e, float h, vec3 wind, float r, int texID) {
    // y positions
    vec3 ups = vec3(0, s*h, 0);
    vec3 upe = vec3(0, e*h, 0);

    // setup all positions
    vec3 p1 = root - right + upe;
    vec3 p2 = root - right + ups;
    vec3 p3 = root + right + upe;
    vec3 p4 = root + right + ups;

    // displace by wind
    float sx = s*s;
    float ex = e*e;
    p1 = displace(p1, root, wind, ex, r);
    p2 = displace(p2, root, wind, sx, r);
    p3 = displace(p3, root, wind, ex, r);
    p4 = displace(p4, root, wind, sx, r);

    // calc normal
    vec3 n = calcNormal(p1, p2, p3);

    // create grass blade segment
    o.position = p1;
    o.uv       = vec2(0, 1-e);
    o.normal   = n;
    o.texID    = texID;
    gl_Position = P*V*M * vec4(p1, 1.0);
    EmitVertex();
    o.position = p2;
    o.uv       = vec2(0, 1-s);
    o.normal   = n;
    o.texID    = texID;
    gl_Position = P*V*M * vec4(p2, 1.0);
    EmitVertex();
    o.position = p3;
    o.uv       = vec2(1, 1-e);
    o.normal   = n;
    o.texID    = texID;
    gl_Position = P*V*M * vec4(p3, 1.0);
    EmitVertex();
    o.position = p4;
    o.uv       = vec2(1, 1-s);
    o.normal   = n;
    o.texID    = texID;
    gl_Position = P*V*M * vec4(p4, 1.0);
    EmitVertex();
    EndPrimitive();
}
void lod3(vec3 root, vec3 right, float h, vec3 wind, float r, int texID) {
    // grass blade consists of 5 segments
    float s = 1/5.0;
    makeSegment(root, right, 0.0, 1*s, h, wind, r, texID);
    makeSegment(root, right, 1*s, 2*s, h, wind, r, texID);
    makeSegment(root, right, 2*s, 3*s, h, wind, r, texID);
    makeSegment(root, right, 3*s, 4*s, h, wind, r, texID);
    makeSegment(root, right, 4*s, 1.0, h, wind, r, texID);
}
void lod2(vec3 root, vec3 right, float h, vec3 wind, float r, int texID) {
    // grass blade consists of 3 segments
    float s = 1/3.0;
    makeSegment(root, right, 0.0, 1*s, h, wind, r, texID);
    makeSegment(root, right, 1*s, 2*s, h, wind, r, texID);
    makeSegment(root, right, 2*s, 1.0, h, wind, r, texID);
}
void lod1(vec3 root, vec3 right, float h, vec3 wind, float r, int texID) {
    // grass blade consists of 1 segment
    makeSegment(root, right, 0.0, 1.0, h, wind, r, texID);
}
void lod0(Tile tile, int texID) {
    // get tile radius in x and z
    vec2 tx = vec2(tilesize/2, 0);
    vec2 tz = vec2(0, tilesize/2);

    // four corners of the tile
    vec2 p1 = tile.pos - tx + tz;
    vec2 p2 = tile.pos - tx - tz;
    vec2 p3 = tile.pos + tx + tz;
    vec2 p4 = tile.pos + tx - tz;
    // calc corresponding heights
    float h1 = calcRootHeight(tile.tri1, p1.x, p1.y).y + 10;
    float h2 = calcRootHeight(tile.tri1, p2.x, p2.y).y + 10;
    float h3 = calcRootHeight(tile.tri1, p3.x, p3.y).y + 10;
    float h4 = calcRootHeight(tile.tri2, p4.x, p4.y).y + 10;
    // vertices
    vec3 v1 = vec3(p1.x, h1, p1.y);
    vec3 v2 = vec3(p2.x, h2, p2.y);
    vec3 v3 = vec3(p3.x, h3, p3.y);
    vec3 v4 = vec3(p4.x, h4, p4.y);

    // calc normal
    vec3 n = calcNormal(v1, v2, v3);

    // emit vertices for tile
    o.position = v1;
    o.uv       = vec2(0, 0);
    o.normal   = n;
    o.texID    = texID;
    gl_Position = P*V*M * vec4(v1, 1.0);
    EmitVertex();
    o.position = v2;
    o.uv       = vec2(0, 1);
    o.normal   = n;
    o.texID    = texID;
    gl_Position = P*V*M * vec4(v2, 1.0);
    EmitVertex();
    o.position = v3;
    o.uv       = vec2(1, 0);
    o.normal   = n;
    o.texID    = texID;
    gl_Position = P*V*M * vec4(v3, 1.0);
    EmitVertex();
    o.position = v4;
    o.uv       = vec2(1, 1);
    o.normal   = n;
    o.texID    = texID;
    gl_Position = P*V*M * vec4(v4, 1.0);
    EmitVertex();
    EndPrimitive();
}

//-----------------------------------------------------------------------------------//
// calculate texture                                                                 //
//-----------------------------------------------------------------------------------//

int getTextureID(float r) {
    int              texID = 3;
    if(r > 0.8)      texID = 1;
    else if(r > 0.5) texID = 2;
    return texID;
}

//-----------------------------------------------------------------------------------//
// calculate wind                                                                    //
//-----------------------------------------------------------------------------------//
vec2 getWindAt(int x, int z) {
    vec2 wind = vec2(0, 0);
    if(abs(x) <= radius && abs(z) <= radius) {
        int dim = 2*radius + 1;
        int idx = (z+radius)*dim + (x+radius);
        vec4 vel = velocity[idx];
        wind = vec2(vel.x, vel.y);
    }
    return wind;
}
vec3 calcWind(Tile tile, vec3 root, vec3 local) {
    // extract tile positions of the current tile and camera
    int tx = int(tile.pos.x / tilesize);
    int tz = int(tile.pos.y / tilesize);
    int cx = int(cameraPos.x / tilesize);
    int cz = int(cameraPos.z / tilesize);
    if(tile.pos.x < 0) { tx -= 1; }
    if(tile.pos.y < 0) { tz -= 1; }
    if(cameraPos.x < 0) { cx -= 1; }
    if(cameraPos.z < 0) { cz -= 1; }

    // tile relative to camera pos
    int rx = tx - cx;
    int rz = tz - cz;

    // get interpolation alpha from relative grass position
    vec3 tilePos = vec3(tile.pos.x, 0.0, tile.pos.y);
    int   dx = int(sign(local.x));
    int   dz = int(sign(local.z));
    float ax = abs(local.x);
    float az = abs(local.z); 

    // make wind by bilear interpolation
    vec2 w0  = getWindAt(rx,    rz   );
    vec2 w0x = getWindAt(rx+dx, rz   );
    vec2 w1  = getWindAt(rx,    rz+dz);
    vec2 w1x = getWindAt(rx+dx, rz+dz);
    w0        = mix(w0, w0x, ax);
    w1        = mix(w1, w1x, ax);
    vec2 wind = mix(w0, w1,  az);

    return vec3(wind.x, 0, wind.y);
}

void main() {
    // current tile
    Tile tile = tiles[i[0].id];
    // current vertex ID
    int vid = i[0].vid;

    // random numbers
    float r = rand();

    // setup vectors
    vec3  local  = getRootLocalPos(r);
    vec3  root   = calcRootWorldPos(local);
    vec3  right  = extractCameraRight()*calcLODBladeWidth(root);
    float height = range(45, 50)*calcLODBladeHeight(root);
    vec3  wind   = calcWind(tile, root, local);

    // calc blade count
    if(vid > calcLODBladeCount(root)) {
        EndPrimitive();
    } else {
        // create segments depending on level of detail
        int lod = calcLOD(root);
        int texID = getTextureID(r);
        if     (lod == 3) { lod3(root, right, height, wind, r, texID); }
        else if(lod == 2) { lod2(root, right, height, wind, r, texID); }
        else if(lod == 1) { lod1(root, right, height, wind, r, texID); }
        else              { lod0(tile, 0);                     }
    }
}