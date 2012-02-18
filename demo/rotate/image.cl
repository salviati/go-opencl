constant sampler_t linear  = CLK_NORMALIZED_COORDS_FALSE | CLK_ADDRESS_CLAMP | CLK_FILTER_LINEAR;
constant sampler_t nearest = CLK_NORMALIZED_COORDS_FALSE | CLK_ADDRESS_CLAMP | CLK_FILTER_NEAREST;

// recfactorx = 1/factorx ---taking a division outside of the "loop".
__kernel void image_recscale(__read_only image2d_t src, __write_only image2d_t dst, float recfactorx, float recfactory) {
	int2 p = {get_global_id(0), get_global_id(1)};
	float2 q = {convert_float(p.x)*recfactorx, convert_float(p.y)*recfactory};

	uint4 pixel = read_imageui(src, linear, q);
	write_imageui(dst, p, pixel);
}

__kernel void image_rotate(__read_only  image2d_t src, __write_only image2d_t dst, float angle) {
	int2 p = {get_global_id(0), get_global_id(1)};
	int2 o = {get_image_width(src)/2, get_image_height(src)/2};

	float2 po = convert_float2(p-o);

	float c;
	float s = sincos(angle, &c);

	float2 qo = {c*po.x - s*po.y, s*po.x + c*po.y};
	float2 q = (qo) + convert_float2(o);

	uint4 pixel = read_imageui(src, linear, q);
	write_imageui(dst, p, pixel);
}

__kernel void image_flip_h(__read_only  image2d_t src, __write_only image2d_t dst) {
	int2 p = {get_global_id(0), get_global_id(1)};
	int2 q = {get_image_width(src) - p.x, p.y};
	
	uint4 pixel = read_imageui(src, nearest, q);
	write_imageui(dst, p, pixel);
}

__kernel void image_flip_v(__read_only  image2d_t src, __write_only image2d_t dst) {
	int2 p = {get_global_id(0), get_global_id(1)};
	int2 q = {p.x, get_image_height(src) - p.y};
	
	uint4 pixel = read_imageui(src, nearest, q);
	write_imageui(dst, p, pixel);
}

__kernel void image_flip_hv(__read_only  image2d_t src, __write_only image2d_t dst) {
	int2 p = {get_global_id(0), get_global_id(1)};
	int2 q = {get_image_width(src) - p.x, get_image_height(src) - p.y};
	
	uint4 pixel = read_imageui(src, nearest, q);
	write_imageui(dst, p, pixel);
}

__kernel void image_affine(__read_only  image2d_t src, __write_only image2d_t dst, float2 Ax, float2 Ay) {
	int2 p = {get_global_id(0), get_global_id(1)};
	float2 pf = convert_float2(p);
	float2 q = {dot(Ax, pf), dot(Ay, pf)};

	uint4 pixel = read_imageui(src, linear, q);
	write_imageui(dst, p, pixel);
}

__kernel void image_affine2(__read_only  image2d_t src, __write_only image2d_t dst, float2 Ax, float2 Ay, float2 p0, float2 q0) {
	int2 p = {get_global_id(0), get_global_id(1)};
	float2 pf = convert_float2(p);
	pf -= p0;
	float2 q = {dot(Ax, pf), dot(Ay, pf)};
	q += q0;

	uint4 pixel = read_imageui(src, linear, q);
	write_imageui(dst, p, pixel);
}
