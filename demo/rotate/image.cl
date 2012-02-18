constant sampler_t sampler = CLK_NORMALIZED_COORDS_FALSE | CLK_ADDRESS_CLAMP | CLK_FILTER_LINEAR;

__kernel void image_scale(read_only image2d_t src, write_only image2d_t dst, float factorx, float factory) {
	int2 p = {get_global_id(0), get_global_id(1)};
	float2 q = {convert_float(get_global_id(0))/factorx, convert_float(get_global_id(1))/factory};

	uint4 pixel = read_imageui(src, sampler, q);
	write_imageui(dst, p, pixel);
}

__kernel void image_rotate(__read_only  image2d_t src, __write_only image2d_t dst, float angle) {
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

__kernel void image_flip_h(__read_only  image2d_t src, __write_only image2d_t dst) {
	int2 p = {get_global_id(0), get_global_id(1)};
	int2 q = {get_image_width(src) - get_global_id(0), get_global_id(1)};
	
	uint4 pixel = read_imageui(src, sampler, q);
	write_imageui(dst, p, pixel);
}

__kernel void image_flip_v(__read_only  image2d_t src, __write_only image2d_t dst) {
	int2 p = {get_global_id(0), get_global_id(1)};
	int2 q = {get_global_id(0), get_image_height(src) - get_global_id(1)};
	
	uint4 pixel = read_imageui(src, sampler, q);
	write_imageui(dst, p, pixel);
}

__kernel void image_flip_hv(__read_only  image2d_t src, __write_only image2d_t dst) {
	int2 p = {get_global_id(0), get_global_id(1)};
	int2 q = {get_image_width(src) - get_global_id(0), get_image_height(src) - get_global_id(1)};
	
	uint4 pixel = read_imageui(src, sampler, q);
	write_imageui(dst, p, pixel);
}
