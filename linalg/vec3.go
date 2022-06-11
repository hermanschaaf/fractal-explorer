package linalg

import "math"

type Vec3 struct {
	X float64
	Y float64
	Z float64
}

func (v1 Vec3) Add(v2 Vec3) Vec3 {
	return Vec3{
		X: v1.X + v2.X,
		Y: v1.Y + v2.Y,
		Z: v1.Z + v2.Z,
	}
}

func (v1 Vec3) Sub(v2 Vec3) Vec3 {
	return Vec3{
		X: v1.X - v2.X,
		Y: v1.Y - v2.Y,
		Z: v1.Z - v2.Z,
	}
}

func (v1 Vec3) Multiply(v float64) Vec3 {
	return Vec3{
		X: v1.X * v,
		Y: v1.Y * v,
		Z: v1.Z * v,
	}
}

func (v1 Vec3) Dot(v2 Vec3) float64 {
	return v1.X*v2.X + v1.Y*v2.Y + v1.Z*v2.Z
}

func (v1 Vec3) Multiply2D(v2 Vec3) Vec3 {
	a, b := v1.X, v1.Y
	c, d := v2.X, v2.Y
	return Vec3{
		X: a*c - b*d, // use a minus here to simulate complex numbers
		Y: a*d + b*c,
		Z: 0,
	}
}

func (v1 Vec3) Multiply3D(v2 Vec3) Vec3 {
	a, b, c := v1.X, v1.Y, v1.Z
	d, e, f := v2.X, v2.Y, v2.Z
	return Vec3{
		X: a*d - b*e + c*f,
		Y: a*e + b*d - a*f,
		Z: c*d + c*f - c*e,
	}
}

func (v1 Vec3) Cross(v2 Vec3) Vec3 {
	return Vec3{
		X: v1.Y*v2.Z - v1.Z*v2.Y,
		Y: v1.Z*v2.X - v1.X*v2.Z,
		Z: v1.X*v2.Y - v1.Y*v2.X,
	}
}

func (v1 Vec3) Length() float64 {
	return math.Sqrt(v1.X*v1.X + v1.Y*v1.Y + v1.Z*v1.Z)
}

func (v1 Vec3) Normalize() Vec3 {
	l := v1.Length()
	if l == 0 {
		return Vec3{
			X: 1,
			Y: 1,
			Z: 1,
		}
	}
	return Vec3{
		X: v1.X / l,
		Y: v1.Y / l,
		Z: v1.Z / l,
	}
}

// RotateX rotates the vector around the Y axis
func (v1 Vec3) RotateX(angle float64) Vec3 {
	return Vec3{
		X: v1.X,
		Y: math.Cos(angle)*v1.Y - math.Sin(angle)*v1.Z,
		Z: math.Sin(angle)*v1.Y + math.Cos(angle)*v1.Z,
	}
}

// RotateY rotates the vector around the Y axis
func (v1 Vec3) RotateY(angle float64) Vec3 {
	return Vec3{
		X: math.Cos(angle)*v1.X + math.Sin(angle)*v1.Z,
		Y: v1.Y,
		Z: -math.Sin(angle)*v1.X + math.Cos(angle)*v1.Z,
	}
}

// RotateZ rotates the vector around the Z axis
func (v1 Vec3) RotateZ(angle float64) Vec3 {
	return Vec3{
		X: math.Cos(angle)*v1.X - math.Sin(angle)*v1.Y,
		Y: math.Sin(angle)*v1.X + math.Cos(angle)*v1.Y,
		Z: v1.Z,
	}
}

// ProperRotation performs a proper rotation R by angle θ around the axis u = (ux, uy, uz),
// a unit vector with u_x^2 + u_y^2 + u_z^2 = x1
func (v1 Vec3) ProperRotation(angle float64, u Vec3) Vec3 {
	cosAngle := math.Cos(angle)
	sinAngle := math.Sin(angle)

	x1 := cosAngle + sq(u.X)*(1-cosAngle)
	x2 := u.X*u.Y*(1-cosAngle) - u.Z*sinAngle
	x3 := u.X*u.Z*(1-cosAngle) + u.Y*sinAngle

	y1 := u.Y*u.X*(1-cosAngle) + u.Z*sinAngle
	y2 := cosAngle + sq(u.Y)*(1-cosAngle)
	y3 := u.Y*u.Z*(1-cosAngle) - u.X*sinAngle

	z1 := u.Z*u.X*(1-cosAngle) - u.Y*sinAngle
	z2 := u.Z*u.Y*(1-cosAngle) + u.X*sinAngle
	z3 := cosAngle + sq(u.Z)*(1-cosAngle)

	return Vec3{
		X: x1*v1.X + x2*v1.Y + x3*v1.Z,
		Y: y1*v1.X + y2*v1.Y + y3*v1.Z,
		Z: z1*v1.X + z2*v1.Y + z3*v1.Z,
	}
}

func sq(v float64) float64 {
	return math.Pow(v, 2)
}
