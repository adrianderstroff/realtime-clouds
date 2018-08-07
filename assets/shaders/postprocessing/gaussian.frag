#version 410

in VertexOut {
    vec2 uv;
}i;
out vec3 fragColor;

//declare uniforms
uniform sampler2D texture;
uniform float     resolution;
uniform float     radius;
uniform vec2      dir;

vec4 tap5(vec2 tc, float blur, float hstep, float vstep) {
	vec4 sum = vec4(0.0);

	sum += texture(texture, vec2(tc.x - 2.0*blur*hstep, tc.y - 2.0*blur*vstep)) * 0.06136;
	sum += texture(texture, vec2(tc.x - 1.0*blur*hstep, tc.y - 1.0*blur*vstep)) * 0.24477;
	
	sum += texture(texture, vec2(tc.x, tc.y)) * 0.38774;
	
	sum += texture(texture, vec2(tc.x + 1.0*blur*hstep, tc.y + 1.0*blur*vstep)) * 0.24477;
	sum += texture(texture, vec2(tc.x + 2.0*blur*hstep, tc.y + 2.0*blur*vstep)) * 0.06136;

	return sum;
}

vec4 tap9(vec2 tc, float blur, float hstep, float vstep) {
	vec4 sum = vec4(0.0);

	sum += texture(texture, vec2(tc.x - 4.0*blur*hstep, tc.y - 4.0*blur*vstep)) * 0.0162162162;
	sum += texture(texture, vec2(tc.x - 3.0*blur*hstep, tc.y - 3.0*blur*vstep)) * 0.0540540541;
	sum += texture(texture, vec2(tc.x - 2.0*blur*hstep, tc.y - 2.0*blur*vstep)) * 0.1216216216;
	sum += texture(texture, vec2(tc.x - 1.0*blur*hstep, tc.y - 1.0*blur*vstep)) * 0.1945945946;
	
	sum += texture(texture, vec2(tc.x, tc.y)) * 0.2270270270;
	
	sum += texture(texture, vec2(tc.x + 1.0*blur*hstep, tc.y + 1.0*blur*vstep)) * 0.1945945946;
	sum += texture(texture, vec2(tc.x + 2.0*blur*hstep, tc.y + 2.0*blur*vstep)) * 0.1216216216;
	sum += texture(texture, vec2(tc.x + 3.0*blur*hstep, tc.y + 3.0*blur*vstep)) * 0.0540540541;
	sum += texture(texture, vec2(tc.x + 4.0*blur*hstep, tc.y + 4.0*blur*vstep)) * 0.0162162162;

	return sum;
}
vec4 tap13(vec2 tc, float blur, float hstep, float vstep) {
	vec4 sum = vec4(0.0);

	sum += texture(texture, vec2(tc.x - 6.0*blur*hstep, tc.y - 6.0*blur*vstep)) * 0.002406;
	sum += texture(texture, vec2(tc.x - 5.0*blur*hstep, tc.y - 5.0*blur*vstep)) * 0.009255;
	sum += texture(texture, vec2(tc.x - 4.0*blur*hstep, tc.y - 4.0*blur*vstep)) * 0.027867;
	sum += texture(texture, vec2(tc.x - 3.0*blur*hstep, tc.y - 3.0*blur*vstep)) * 0.065666;
	sum += texture(texture, vec2(tc.x - 2.0*blur*hstep, tc.y - 2.0*blur*vstep)) * 0.121117;
	sum += texture(texture, vec2(tc.x - 1.0*blur*hstep, tc.y - 1.0*blur*vstep)) * 0.174868;
	
	sum += texture(texture, vec2(tc.x, tc.y)) * 0.197641;
	
	sum += texture(texture, vec2(tc.x + 1.0*blur*hstep, tc.y + 1.0*blur*vstep)) * 0.174868;
	sum += texture(texture, vec2(tc.x + 2.0*blur*hstep, tc.y + 2.0*blur*vstep)) * 0.121117;
	sum += texture(texture, vec2(tc.x + 3.0*blur*hstep, tc.y + 3.0*blur*vstep)) * 0.065666;
	sum += texture(texture, vec2(tc.x + 4.0*blur*hstep, tc.y + 4.0*blur*vstep)) * 0.027867;
	sum += texture(texture, vec2(tc.x + 5.0*blur*hstep, tc.y + 5.0*blur*vstep)) * 0.009255;
	sum += texture(texture, vec2(tc.x + 6.0*blur*hstep, tc.y + 6.0*blur*vstep)) * 0.002406;

	return sum;
}

void main() {
	//our original texcoord for this fragment
	vec2 tc = i.uv;
	
	//the amount to blur, i.e. how far off center to sample from 
	float blur = radius/resolution; 
    
	//the direction of our blur
	float hstep = dir.x;
	float vstep = dir.y;
    
	//apply blurring, using a 9-tap filter with predefined gaussian weights
    vec4 sum = tap13(tc, blur, hstep, vstep);
	
	fragColor =  sum.rgb;
}