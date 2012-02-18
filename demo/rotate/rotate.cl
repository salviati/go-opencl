const sampler_t sampler  = CLK_NORMALIZED_COORDS_FALSE | CLK_ADDRESS_CLAMP | CLK_FILTER_LINEAR;


// p: point in dst
// q: point in src
// o: rotate about this point
// po: point in dst, relative to o
// qo: point in dst, relative to o
// R = {Rx,Ry} is the rotation matrix

__kernel void rotateImage(__read_only  image2d_t src, __write_only image2d_t dst, float angle)
{
	int2 p = {get_global_id(0), get_global_id(1)};
	int2 o = {get_image_width(src)/2, get_image_height(src)/2};

	float2 po = convert_float2(p-o);

	float c;
	float s = sincos(angle, &c);
	float2 Rx = {c,-s};
	float2 Ry = {s,c};
	
	float2 qo = {dot(Rx, po), dot(Ry, po)};
	float2 q = (qo) + convert_float2(o);

	uint4 pixel = read_imageui(src, sampler, q);
	write_imageui(dst, p, pixel);
}
