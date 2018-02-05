// Copyright ©2018 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This is a translation of the FFTPACK rfft functions by
// Paul N Swarztrauber, placed in the public domain at
// http://www.netlib.org/fftpack/.

package fourier

import "math"

// cffti initializes the array work which is used in both cfftf
// and cfftb. the prime factorization of n together with a
// tabulation of the trigonometric functions are computed and
// stored in work.
//
//  input parameter
//
//  n      the length of the sequence to be transformed.
//
//  output parameters
//
//  work   a work array which must be dimensioned at least 4*n.
//         the same work array can be used for both cfftf and cfftb
//         as long as n remains unchanged. different work arrays
//         are required for different values of n. the contents of
//         work must not be changed between calls of cfftf or cfftb.
//
//  ifac   a work array containing the factors of n. ifac must have
//         length 15.
func cffti(n int, work []float64, ifac []int) {
	if n == 1 {
		return
	}
	cffti1(n, work[2*n:], ifac)
}

func cffti1(n int, wa []float64, ifac []int) {
	ntryh := [4]int{3, 4, 2, 5}

	nl := n
	nf := 0

outer:
	for j, ntry := 0, 0; ; j++ {
		if j < 4 {
			ntry = ntryh[j]
		} else {
			ntry += 2
		}
		for {
			if nl%ntry != 0 {
				continue outer
			}

			ifac[nf+2] = ntry
			nl /= ntry
			nf++

			if ntry == 2 && nf != 1 {
				for i := 1; i < nf; i++ {
					ib := nf - i + 1
					ifac[ib+1] = ifac[ib]
				}
				ifac[2] = 2
			}

			if nl == 1 {
				break outer
			}
		}
	}

	ifac[0] = n
	ifac[1] = nf

	argh := 2 * math.Pi / float64(n)
	i := 1
	l1 := 1
	for k1 := 0; k1 < nf; k1++ {
		ip := ifac[k1+2]
		ld := 0
		l2 := l1 * ip
		ido := n / l2
		idot := 2*ido + 2
		for j := 0; j < ip-1; j++ {
			i1 := i
			wa[i-1] = 1
			wa[i] = 0
			ld += l1
			var fi float64
			argld := float64(ld) * argh
			for ii := 3; ii < idot; ii += 2 {
				i += 2
				fi++
				arg := fi * argld
				wa[i-1] = math.Cos(arg)
				wa[i] = math.Sin(arg)
			}
			if ip > 5 {
				wa[i1-1] = wa[i-1]
				wa[i1] = wa[i]
			}
		}
		l1 = l2
	}
}

// cfftf computes the forward complex discrete fourier
// transform (the fourier analysis). equivalently , cfftf computes
// the fourier coefficients of a complex periodic sequence.
// the transform is defined below at output parameter c.
//
//  input parameters
//
//  n      the length of the array c to be transformed.  the method
//         is most efficient when n is a product of small primes.
//         n may change so long as different work arrays are provided
//
//  c      a complex array of length n which contains the sequence
//         to be transformed
//
//  work   a real work array which must be dimensioned at least 4*n.
//         in the program that calls cfftf. the work array must be
//         initialized by calling subroutine cffti(n,work) and a
//         different work array must be used for each different
//         value of n. this initialization does not have to be
//         repeated so long as n remains unchanged thus subsequent
//         transforms can be obtained faster than the first.
//         the same work array can be used by cfftf and cfftb.
//
//  ifac   a work array containing the factors of n. ifac must have
//         length 15.
//
//  output parameters
//
//   c     for j=1,...,n
//           c(j)=the sum from k=1,...,n of
//             c(k)*exp(-i*(j-1)*(k-1)*2*pi/n)
//
//         where i=sqrt(-1)
//
//  This transform is unnormalized since a call of cfftf
//  followed by a call of cfftb will multiply the input
//  sequence by n.
//
//  work   contains results which must not be destroyed between
//         calls of cfftf or cfftb.
//  ifac   contains results which must not be destroyed between
//         calls of cfftf or cfftb.
func cfftf(n int, r, work []float64, ifac []int) {
	if n == 1 {
		return
	}
	cfftf1(n, r, work, work[2*n:], ifac)
}

func cfftf1(n int, c, ch []float64, wa oneArray, ifac oneIntArray) {
	nf := ifac.at(2)
	na := 0
	l1 := 1
	iw := 1

	for k1 := 1; k1 <= nf; k1++ {
		ip := ifac.at(k1 + 2)
		l2 := ip * l1
		ido := n / l2
		idot := 2 * ido
		idl1 := idot * l1

		switch ip {
		case 4:
			ix2 := iw + idot
			ix3 := ix2 + idot
			if na == 0 {
				passf4(idot, l1, c, ch, wa.sliceFrom(iw), wa.sliceFrom(ix2), wa.sliceFrom(ix3))
			} else {
				passf4(idot, l1, ch, c, wa.sliceFrom(iw), wa.sliceFrom(ix2), wa.sliceFrom(ix3))
			}
			na = 1 - na
		case 2:
			if na == 0 {
				passf2(idot, l1, c, ch, wa.sliceFrom(iw))
			} else {
				passf2(idot, l1, ch, c, wa.sliceFrom(iw))
			}
			na = 1 - na
		case 3:
			ix2 := iw + idot
			if na == 0 {
				passf3(idot, l1, c, ch, wa.sliceFrom(iw), wa.sliceFrom(ix2))
			} else {
				passf3(idot, l1, ch, c, wa.sliceFrom(iw), wa.sliceFrom(ix2))
			}
			na = 1 - na
		case 5:
			ix2 := iw + idot
			ix3 := ix2 + idot
			ix4 := ix3 + idot
			if na == 0 {
				passf5(idot, l1, c, ch, wa.sliceFrom(iw), wa.sliceFrom(ix2), wa.sliceFrom(ix3), wa.sliceFrom(ix4))
			} else {
				passf5(idot, l1, ch, c, wa.sliceFrom(iw), wa.sliceFrom(ix2), wa.sliceFrom(ix3), wa.sliceFrom(ix4))
			}
			na = 1 - na
		default:
			var nac bool
			if na == 0 {
				nac = passf(idot, ip, l1, idl1, c, c, c, ch, ch, wa.sliceFrom(iw))
			} else {
				nac = passf(idot, ip, l1, idl1, ch, ch, ch, c, c, wa.sliceFrom(iw))
			}
			if nac {
				na = 1 - na
			}
		}

		l1 = l2
		iw += (ip - 1) * idot
	}

	if na == 0 {
		return
	}
	for i := 0; i < 2*n; i++ {
		c[i] = ch[i]
	}
}

func passf2(ido, l1 int, cc, ch []float64, wa1 oneArray) {
	cc3 := newThreeArray(ido, 2, l1, cc)
	ch3 := newThreeArray(ido, l1, 2, ch)

	if ido <= 2 {
		for k := 1; k <= l1; k++ {
			ch3.set(1, k, 1, cc3.at(1, 1, k)+cc3.at(1, 2, k))
			ch3.set(1, k, 2, cc3.at(1, 1, k)-cc3.at(1, 2, k))
			ch3.set(2, k, 1, cc3.at(2, 1, k)+cc3.at(2, 2, k))
			ch3.set(2, k, 2, cc3.at(2, 1, k)-cc3.at(2, 2, k))
		}
		return
	}
	for k := 1; k <= l1; k++ {
		for i := 2; i <= ido; i += 2 {
			ch3.set(i-1, k, 1, cc3.at(i-1, 1, k)+cc3.at(i-1, 2, k))
			tr2 := cc3.at(i-1, 1, k) - cc3.at(i-1, 2, k)
			ch3.set(i, k, 1, cc3.at(i, 1, k)+cc3.at(i, 2, k))
			ti2 := cc3.at(i, 1, k) - cc3.at(i, 2, k)
			ch3.set(i, k, 2, wa1.at(i-1)*ti2-wa1.at(i)*tr2)
			ch3.set(i-1, k, 2, wa1.at(i-1)*tr2+wa1.at(i)*ti2)
		}
	}
}

func passf3(ido, l1 int, cc, ch []float64, wa1, wa2 oneArray) {
	const (
		taur = -0.5
		taui = -0.866025403784439 // -sqrt(3)/2
	)

	cc3 := newThreeArray(ido, 3, l1, cc)
	ch3 := newThreeArray(ido, l1, 3, ch)

	if ido == 2 {
		for k := 1; k <= l1; k++ {
			tr2 := cc3.at(1, 2, k) + cc3.at(1, 3, k)
			cr2 := cc3.at(1, 1, k) + taur*tr2
			ch3.set(1, k, 1, cc3.at(1, 1, k)+tr2)
			ti2 := cc3.at(2, 2, k) + cc3.at(2, 3, k)
			ci2 := cc3.at(2, 1, k) + taur*ti2
			ch3.set(2, k, 1, cc3.at(2, 1, k)+ti2)
			cr3 := taui * (cc3.at(1, 2, k) - cc3.at(1, 3, k))
			ci3 := taui * (cc3.at(2, 2, k) - cc3.at(2, 3, k))
			ch3.set(1, k, 2, cr2-ci3)
			ch3.set(1, k, 3, cr2+ci3)
			ch3.set(2, k, 2, ci2+cr3)
			ch3.set(2, k, 3, ci2-cr3)
		}
		return
	}
	for k := 1; k <= l1; k++ {
		for i := 2; i <= ido; i += 2 {
			tr2 := cc3.at(i-1, 2, k) + cc3.at(i-1, 3, k)
			cr2 := cc3.at(i-1, 1, k) + taur*tr2
			ch3.set(i-1, k, 1, cc3.at(i-1, 1, k)+tr2)
			ti2 := cc3.at(i, 2, k) + cc3.at(i, 3, k)
			ci2 := cc3.at(i, 1, k) + taur*ti2
			ch3.set(i, k, 1, cc3.at(i, 1, k)+ti2)
			cr3 := taui * (cc3.at(i-1, 2, k) - cc3.at(i-1, 3, k))
			ci3 := taui * (cc3.at(i, 2, k) - cc3.at(i, 3, k))
			dr2 := cr2 - ci3
			dr3 := cr2 + ci3
			di2 := ci2 + cr3
			di3 := ci2 - cr3
			ch3.set(i, k, 2, wa1.at(i-1)*di2-wa1.at(i)*dr2)
			ch3.set(i-1, k, 2, wa1.at(i-1)*dr2+wa1.at(i)*di2)
			ch3.set(i, k, 3, wa2.at(i-1)*di3-wa2.at(i)*dr3)
			ch3.set(i-1, k, 3, wa2.at(i-1)*dr3+wa2.at(i)*di3)
		}
	}
}

func passf4(ido, l1 int, cc, ch []float64, wa1, wa2, wa3 oneArray) {
	cc3 := newThreeArray(ido, 4, l1, cc)
	ch3 := newThreeArray(ido, l1, 4, ch)

	if ido == 2 {
		for k := 1; k <= l1; k++ {
			ti1 := cc3.at(2, 1, k) - cc3.at(2, 3, k)
			ti2 := cc3.at(2, 1, k) + cc3.at(2, 3, k)
			tr4 := cc3.at(2, 2, k) - cc3.at(2, 4, k)
			ti3 := cc3.at(2, 2, k) + cc3.at(2, 4, k)
			tr1 := cc3.at(1, 1, k) - cc3.at(1, 3, k)
			tr2 := cc3.at(1, 1, k) + cc3.at(1, 3, k)
			ti4 := cc3.at(1, 4, k) - cc3.at(1, 2, k)
			tr3 := cc3.at(1, 2, k) + cc3.at(1, 4, k)
			ch3.set(1, k, 1, tr2+tr3)
			ch3.set(1, k, 3, tr2-tr3)
			ch3.set(2, k, 1, ti2+ti3)
			ch3.set(2, k, 3, ti2-ti3)
			ch3.set(1, k, 2, tr1+tr4)
			ch3.set(1, k, 4, tr1-tr4)
			ch3.set(2, k, 2, ti1+ti4)
			ch3.set(2, k, 4, ti1-ti4)
		}
		return
	}
	for k := 1; k <= l1; k++ {
		for i := 2; i <= ido; i += 2 {
			ti1 := cc3.at(i, 1, k) - cc3.at(i, 3, k)
			ti2 := cc3.at(i, 1, k) + cc3.at(i, 3, k)
			ti3 := cc3.at(i, 2, k) + cc3.at(i, 4, k)
			tr4 := cc3.at(i, 2, k) - cc3.at(i, 4, k)
			tr1 := cc3.at(i-1, 1, k) - cc3.at(i-1, 3, k)
			tr2 := cc3.at(i-1, 1, k) + cc3.at(i-1, 3, k)
			ti4 := cc3.at(i-1, 4, k) - cc3.at(i-1, 2, k)
			tr3 := cc3.at(i-1, 2, k) + cc3.at(i-1, 4, k)
			ch3.set(i-1, k, 1, tr2+tr3)
			cr3 := tr2 - tr3
			ch3.set(i, k, 1, ti2+ti3)
			ci3 := ti2 - ti3
			cr2 := tr1 + tr4
			cr4 := tr1 - tr4
			ci2 := ti1 + ti4
			ci4 := ti1 - ti4
			ch3.set(i-1, k, 2, wa1.at(i-1)*cr2+wa1.at(i)*ci2)
			ch3.set(i, k, 2, wa1.at(i-1)*ci2-wa1.at(i)*cr2)
			ch3.set(i-1, k, 3, wa2.at(i-1)*cr3+wa2.at(i)*ci3)
			ch3.set(i, k, 3, wa2.at(i-1)*ci3-wa2.at(i)*cr3)
			ch3.set(i-1, k, 4, wa3.at(i-1)*cr4+wa3.at(i)*ci4)
			ch3.set(i, k, 4, wa3.at(i-1)*ci4-wa3.at(i)*cr4)
		}
	}
}

func passf5(ido, l1 int, cc, ch []float64, wa1, wa2, wa3, wa4 oneArray) {
	const (
		tr11 = 0.309016994374947
		ti11 = -0.951056516295154
		tr12 = -0.809016994374947
		ti12 = -0.587785252292473
	)

	cc3 := newThreeArray(ido, 5, l1, cc)
	ch3 := newThreeArray(ido, l1, 5, ch)

	if ido == 2 {
		for k := 1; k <= l1; k++ {
			ti5 := cc3.at(2, 2, k) - cc3.at(2, 5, k)
			ti2 := cc3.at(2, 2, k) + cc3.at(2, 5, k)
			ti4 := cc3.at(2, 3, k) - cc3.at(2, 4, k)
			ti3 := cc3.at(2, 3, k) + cc3.at(2, 4, k)
			tr5 := cc3.at(1, 2, k) - cc3.at(1, 5, k)
			tr2 := cc3.at(1, 2, k) + cc3.at(1, 5, k)
			tr4 := cc3.at(1, 3, k) - cc3.at(1, 4, k)
			tr3 := cc3.at(1, 3, k) + cc3.at(1, 4, k)
			ch3.set(1, k, 1, cc3.at(1, 1, k)+tr2+tr3)
			ch3.set(2, k, 1, cc3.at(2, 1, k)+ti2+ti3)
			cr2 := cc3.at(1, 1, k) + tr11*tr2 + tr12*tr3
			ci2 := cc3.at(2, 1, k) + tr11*ti2 + tr12*ti3
			cr3 := cc3.at(1, 1, k) + tr12*tr2 + tr11*tr3
			ci3 := cc3.at(2, 1, k) + tr12*ti2 + tr11*ti3
			cr5 := ti11*tr5 + ti12*tr4
			ci5 := ti11*ti5 + ti12*ti4
			cr4 := ti12*tr5 - ti11*tr4
			ci4 := ti12*ti5 - ti11*ti4
			ch3.set(1, k, 2, cr2-ci5)
			ch3.set(1, k, 5, cr2+ci5)
			ch3.set(2, k, 2, ci2+cr5)
			ch3.set(2, k, 3, ci3+cr4)
			ch3.set(1, k, 3, cr3-ci4)
			ch3.set(1, k, 4, cr3+ci4)
			ch3.set(2, k, 4, ci3-cr4)
			ch3.set(2, k, 5, ci2-cr5)
		}
		return
	}
	for k := 1; k <= l1; k++ {
		for i := 2; i <= ido; i += 2 {
			ti5 := cc3.at(i, 2, k) - cc3.at(i, 5, k)
			ti2 := cc3.at(i, 2, k) + cc3.at(i, 5, k)
			ti4 := cc3.at(i, 3, k) - cc3.at(i, 4, k)
			ti3 := cc3.at(i, 3, k) + cc3.at(i, 4, k)
			tr5 := cc3.at(i-1, 2, k) - cc3.at(i-1, 5, k)
			tr2 := cc3.at(i-1, 2, k) + cc3.at(i-1, 5, k)
			tr4 := cc3.at(i-1, 3, k) - cc3.at(i-1, 4, k)
			tr3 := cc3.at(i-1, 3, k) + cc3.at(i-1, 4, k)
			ch3.set(i-1, k, 1, cc3.at(i-1, 1, k)+tr2+tr3)
			ch3.set(i, k, 1, cc3.at(i, 1, k)+ti2+ti3)
			cr2 := cc3.at(i-1, 1, k) + tr11*tr2 + tr12*tr3
			ci2 := cc3.at(i, 1, k) + tr11*ti2 + tr12*ti3
			cr3 := cc3.at(i-1, 1, k) + tr12*tr2 + tr11*tr3
			ci3 := cc3.at(i, 1, k) + tr12*ti2 + tr11*ti3
			cr5 := ti11*tr5 + ti12*tr4
			ci5 := ti11*ti5 + ti12*ti4
			cr4 := ti12*tr5 - ti11*tr4
			ci4 := ti12*ti5 - ti11*ti4
			dr3 := cr3 - ci4
			dr4 := cr3 + ci4
			di3 := ci3 + cr4
			di4 := ci3 - cr4
			dr5 := cr2 + ci5
			dr2 := cr2 - ci5
			di5 := ci2 - cr5
			di2 := ci2 + cr5
			ch3.set(i-1, k, 2, wa1.at(i-1)*dr2+wa1.at(i)*di2)
			ch3.set(i, k, 2, wa1.at(i-1)*di2-wa1.at(i)*dr2)
			ch3.set(i-1, k, 3, wa2.at(i-1)*dr3+wa2.at(i)*di3)
			ch3.set(i, k, 3, wa2.at(i-1)*di3-wa2.at(i)*dr3)
			ch3.set(i-1, k, 4, wa3.at(i-1)*dr4+wa3.at(i)*di4)
			ch3.set(i, k, 4, wa3.at(i-1)*di4-wa3.at(i)*dr4)
			ch3.set(i-1, k, 5, wa4.at(i-1)*dr5+wa4.at(i)*di5)
			ch3.set(i, k, 5, wa4.at(i-1)*di5-wa4.at(i)*dr5)
		}
	}
}

func passf(ido, ip, l1, idl1 int, cc, c1, c2, ch, ch2 []float64, wa oneArray) (nac bool) {
	cc3 := newThreeArray(ido, ip, l1, cc)
	c13 := newThreeArray(ido, l1, ip, c1)
	ch3 := newThreeArray(ido, l1, ip, ch)
	c2m := newTwoArray(idl1, ip, c2)
	ch2m := newTwoArray(idl1, ip, ch2)

	idot := ido / 2
	ipp2 := ip + 2
	ipph := (ip + 1) / 2
	idp := ip * ido

	if ido < l1 {
		for j := 2; j <= ipph; j++ {
			jc := ipp2 - j
			for i := 1; i <= ido; i++ {
				for k := 1; k <= l1; k++ {
					ch3.set(i, k, j, cc3.at(i, j, k)+cc3.at(i, jc, k))
					ch3.set(i, k, jc, cc3.at(i, j, k)-cc3.at(i, jc, k))
				}
			}
		}
		for i := 1; i <= ido; i++ {
			for k := 1; k <= l1; k++ {
				ch3.set(i, k, 1, cc3.at(i, 1, k))
			}
		}
	} else {
		for j := 2; j <= ipph; j++ {
			jc := ipp2 - j
			for k := 1; k <= l1; k++ {
				for i := 1; i <= ido; i++ {
					ch3.set(i, k, j, cc3.at(i, j, k)+cc3.at(i, jc, k))
					ch3.set(i, k, jc, cc3.at(i, j, k)-cc3.at(i, jc, k))
				}
			}
		}
		for k := 1; k <= l1; k++ {
			for i := 1; i <= ido; i++ {
				ch3.set(i, k, 1, cc3.at(i, 1, k))
			}
		}
	}

	idl := 2 - ido
	inc := 0
	for l := 2; l <= ipph; l++ {
		lc := ipp2 - l
		idl += ido
		for ik := 1; ik <= idl1; ik++ {
			c2m.set(ik, l, ch2m.at(ik, 1)+wa.at(idl-1)*ch2m.at(ik, 2))
			c2m.set(ik, lc, -wa.at(idl)*ch2m.at(ik, ip))
		}
		idlj := idl
		inc += ido
		for j := 3; j <= ipph; j++ {
			jc := ipp2 - j
			idlj += inc
			if idlj > idp {
				idlj -= idp
			}
			war := wa.at(idlj - 1)
			wai := wa.at(idlj)
			for ik := 1; ik <= idl1; ik++ {
				c2m.set(ik, l, c2m.at(ik, l)+war*ch2m.at(ik, j))
				c2m.set(ik, lc, c2m.at(ik, lc)-wai*ch2m.at(ik, jc))
			}
		}
	}

	for j := 2; j <= ipph; j++ {
		for ik := 1; ik <= idl1; ik++ {
			ch2m.set(ik, 1, ch2m.at(ik, 1)+ch2m.at(ik, j))
		}
	}

	for j := 2; j <= ipph; j++ {
		jc := ipp2 - j
		for ik := 2; ik <= idl1; ik += 2 {
			ch2m.set(ik-1, j, c2m.at(ik-1, j)-c2m.at(ik, jc))
			ch2m.set(ik-1, jc, c2m.at(ik-1, j)+c2m.at(ik, jc))
			ch2m.set(ik, j, c2m.at(ik, j)+c2m.at(ik-1, jc))
			ch2m.set(ik, jc, c2m.at(ik, j)-c2m.at(ik-1, jc))
		}
	}

	if ido == 2 {
		return true
	}

	for ik := 1; ik <= idl1; ik++ {
		c2m.set(ik, 1, ch2m.at(ik, 1))
	}

	for j := 2; j <= ip; j++ {
		for k := 1; k <= l1; k++ {
			c13.set(1, k, j, ch3.at(1, k, j))
			c13.set(2, k, j, ch3.at(2, k, j))
		}
	}

	if idot > l1 {
		idj := 2 - ido
		for j := 2; j <= ip; j++ {
			idj += ido
			for k := 1; k <= l1; k++ {
				idij := idj
				for i := 4; i <= ido; i += 2 {
					idij += 2
					c13.set(i-1, k, j, wa.at(idij-1)*ch3.at(i-1, k, j)+wa.at(idij)*ch3.at(i, k, j))
					c13.set(i, k, j, wa.at(idij-1)*ch3.at(i, k, j)-wa.at(idij)*ch3.at(i-1, k, j))
				}
			}
		}
		return false
	}

	idij := 0
	for j := 2; j <= ip; j++ {
		idij += 2
		for i := 4; i <= ido; i += 2 {
			idij += 2
			for k := 1; k <= l1; k++ {
				c13.set(i-1, k, j, wa.at(idij-1)*ch3.at(i-1, k, j)+wa.at(idij)*ch3.at(i, k, j))
				c13.set(i, k, j, wa.at(idij-1)*ch3.at(i, k, j)-wa.at(idij)*ch3.at(i-1, k, j))
			}
		}
	}
	return false
}

// TODO(kortschak): As described in the documentation below, the
// only change between cfftf and cfftb is the sign before the i in
// the inner loop. The entirety of the code below cfft? could be
// shared between cfftf and cfftb by passing a sign parameter.

// cfftb computes the forward complex discrete fourier
// transform (the fourier analysis). equivalently , cfftb computes
// the fourier coefficients of a complex periodic sequence.
// the transform is defined below at output parameter c.
//
//  input parameters
//
//  n      the length of the array c to be transformed.  the method
//         is most efficient when n is a product of small primes.
//         n may change so long as different work arrays are provided
//
//  c      a complex array of length n which contains the sequence
//         to be transformed
//
//  work   a real work array which must be dimensioned at least 4*n.
//         in the program that calls cfftb. the work array must be
//         initialized by calling subroutine cffti(n,work) and a
//         different work array must be used for each different
//         value of n. this initialization does not have to be
//         repeated so long as n remains unchanged thus subsequent
//         transforms can be obtained faster than the first.
//         the same work array can be used by cfftb and cfftb.
//
//  ifac   a work array containing the factors of n. ifac must have
//         length 15.
//
//  output parameters
//
//   c     for j=1,...,n
//           c(j)=the sum from k=1,...,n of
//             c(k)*exp(i*(j-1)*(k-1)*2*pi/n)
//
//         where i=sqrt(-1)
//
//  This transform is unnormalized since a call of cfftf
//  followed by a call of cfftb will multiply the input
//  sequence by n.
//
//  work   contains results which must not be destroyed between
//         calls of cfftf or cfftb.
//  ifac   contains results which must not be destroyed between
//         calls of cfftf or cfftb.
func cfftb(n int, r, work []float64, ifac []int) {
	if n == 1 {
		return
	}
	cfftb1(n, r, work, work[2*n:], ifac)
}

func cfftb1(n int, c, ch []float64, wa oneArray, ifac oneIntArray) {
	nf := ifac.at(2)
	na := 0
	l1 := 1
	iw := 1

	for k1 := 1; k1 <= nf; k1++ {
		ip := ifac.at(k1 + 2)
		l2 := ip * l1
		ido := n / l2
		idot := 2 * ido
		idl1 := idot * l1

		switch ip {
		case 4:
			ix2 := iw + idot
			ix3 := ix2 + idot
			if na == 0 {
				passb4(idot, l1, c, ch, wa.sliceFrom(iw), wa.sliceFrom(ix2), wa.sliceFrom(ix3))
			} else {
				passb4(idot, l1, ch, c, wa.sliceFrom(iw), wa.sliceFrom(ix2), wa.sliceFrom(ix3))
			}
			na = 1 - na
		case 2:
			if na == 0 {
				passb2(idot, l1, c, ch, wa.sliceFrom(iw))
			} else {
				passb2(idot, l1, ch, c, wa.sliceFrom(iw))
			}
			na = 1 - na
		case 3:
			ix2 := iw + idot
			if na == 0 {
				passb3(idot, l1, c, ch, wa.sliceFrom(iw), wa.sliceFrom(ix2))
			} else {
				passb3(idot, l1, ch, c, wa.sliceFrom(iw), wa.sliceFrom(ix2))
			}
			na = 1 - na
		case 5:
			ix2 := iw + idot
			ix3 := ix2 + idot
			ix4 := ix3 + idot
			if na == 0 {
				passb5(idot, l1, c, ch, wa.sliceFrom(iw), wa.sliceFrom(ix2), wa.sliceFrom(ix3), wa.sliceFrom(ix4))
			} else {
				passb5(idot, l1, ch, c, wa.sliceFrom(iw), wa.sliceFrom(ix2), wa.sliceFrom(ix3), wa.sliceFrom(ix4))
			}
			na = 1 - na
		default:
			var nac bool
			if na == 0 {
				nac = passb(idot, ip, l1, idl1, c, c, c, ch, ch, wa.sliceFrom(iw))
			} else {
				nac = passb(idot, ip, l1, idl1, ch, ch, ch, c, c, wa.sliceFrom(iw))
			}
			if nac {
				na = 1 - na
			}
		}

		l1 = l2
		iw += (ip - 1) * idot
	}

	if na == 0 {
		return
	}
	for i := 0; i < 2*n; i++ {
		c[i] = ch[i]
	}
}

func passb2(ido, l1 int, cc, ch []float64, wa1 oneArray) {
	cc3 := newThreeArray(ido, 2, l1, cc)
	ch3 := newThreeArray(ido, l1, 2, ch)

	if ido <= 2 {
		for k := 1; k <= l1; k++ {
			ch3.set(1, k, 1, cc3.at(1, 1, k)+cc3.at(1, 2, k))
			ch3.set(1, k, 2, cc3.at(1, 1, k)-cc3.at(1, 2, k))
			ch3.set(2, k, 1, cc3.at(2, 1, k)+cc3.at(2, 2, k))
			ch3.set(2, k, 2, cc3.at(2, 1, k)-cc3.at(2, 2, k))
		}
		return
	}
	for k := 1; k <= l1; k++ {
		for i := 2; i <= ido; i += 2 {
			ch3.set(i-1, k, 1, cc3.at(i-1, 1, k)+cc3.at(i-1, 2, k))
			tr2 := cc3.at(i-1, 1, k) - cc3.at(i-1, 2, k)
			ch3.set(i, k, 1, cc3.at(i, 1, k)+cc3.at(i, 2, k))
			ti2 := cc3.at(i, 1, k) - cc3.at(i, 2, k)
			ch3.set(i, k, 2, wa1.at(i-1)*ti2+wa1.at(i)*tr2)
			ch3.set(i-1, k, 2, wa1.at(i-1)*tr2-wa1.at(i)*ti2)
		}
	}
}

func passb3(ido, l1 int, cc, ch []float64, wa1, wa2 oneArray) {
	const (
		taur = -0.5
		taui = 0.866025403784439 // sqrt(3)/2
	)

	cc3 := newThreeArray(ido, 3, l1, cc)
	ch3 := newThreeArray(ido, l1, 3, ch)

	if ido == 2 {
		for k := 1; k <= l1; k++ {
			tr2 := cc3.at(1, 2, k) + cc3.at(1, 3, k)
			cr2 := cc3.at(1, 1, k) + taur*tr2
			ch3.set(1, k, 1, cc3.at(1, 1, k)+tr2)
			ti2 := cc3.at(2, 2, k) + cc3.at(2, 3, k)
			ci2 := cc3.at(2, 1, k) + taur*ti2
			ch3.set(2, k, 1, cc3.at(2, 1, k)+ti2)
			cr3 := taui * (cc3.at(1, 2, k) - cc3.at(1, 3, k))
			ci3 := taui * (cc3.at(2, 2, k) - cc3.at(2, 3, k))
			ch3.set(1, k, 2, cr2-ci3)
			ch3.set(1, k, 3, cr2+ci3)
			ch3.set(2, k, 2, ci2+cr3)
			ch3.set(2, k, 3, ci2-cr3)
		}
		return
	}
	for k := 1; k <= l1; k++ {
		for i := 2; i <= ido; i += 2 {
			tr2 := cc3.at(i-1, 2, k) + cc3.at(i-1, 3, k)
			cr2 := cc3.at(i-1, 1, k) + taur*tr2
			ch3.set(i-1, k, 1, cc3.at(i-1, 1, k)+tr2)
			ti2 := cc3.at(i, 2, k) + cc3.at(i, 3, k)
			ci2 := cc3.at(i, 1, k) + taur*ti2
			ch3.set(i, k, 1, cc3.at(i, 1, k)+ti2)
			cr3 := taui * (cc3.at(i-1, 2, k) - cc3.at(i-1, 3, k))
			ci3 := taui * (cc3.at(i, 2, k) - cc3.at(i, 3, k))
			dr2 := cr2 - ci3
			dr3 := cr2 + ci3
			di2 := ci2 + cr3
			di3 := ci2 - cr3
			ch3.set(i, k, 2, wa1.at(i-1)*di2+wa1.at(i)*dr2)
			ch3.set(i-1, k, 2, wa1.at(i-1)*dr2-wa1.at(i)*di2)
			ch3.set(i, k, 3, wa2.at(i-1)*di3+wa2.at(i)*dr3)
			ch3.set(i-1, k, 3, wa2.at(i-1)*dr3-wa2.at(i)*di3)
		}
	}
}

func passb4(ido, l1 int, cc, ch []float64, wa1, wa2, wa3 oneArray) {
	cc3 := newThreeArray(ido, 4, l1, cc)
	ch3 := newThreeArray(ido, l1, 4, ch)

	if ido == 2 {
		for k := 1; k <= l1; k++ {
			ti1 := cc3.at(2, 1, k) - cc3.at(2, 3, k)
			ti2 := cc3.at(2, 1, k) + cc3.at(2, 3, k)
			tr4 := cc3.at(2, 4, k) - cc3.at(2, 2, k)
			ti3 := cc3.at(2, 2, k) + cc3.at(2, 4, k)
			tr1 := cc3.at(1, 1, k) - cc3.at(1, 3, k)
			tr2 := cc3.at(1, 1, k) + cc3.at(1, 3, k)
			ti4 := cc3.at(1, 2, k) - cc3.at(1, 4, k)
			tr3 := cc3.at(1, 2, k) + cc3.at(1, 4, k)
			ch3.set(1, k, 1, tr2+tr3)
			ch3.set(1, k, 3, tr2-tr3)
			ch3.set(2, k, 1, ti2+ti3)
			ch3.set(2, k, 3, ti2-ti3)
			ch3.set(1, k, 2, tr1+tr4)
			ch3.set(1, k, 4, tr1-tr4)
			ch3.set(2, k, 2, ti1+ti4)
			ch3.set(2, k, 4, ti1-ti4)
		}
		return
	}
	for k := 1; k <= l1; k++ {
		for i := 2; i <= ido; i += 2 {
			ti1 := cc3.at(i, 1, k) - cc3.at(i, 3, k)
			ti2 := cc3.at(i, 1, k) + cc3.at(i, 3, k)
			ti3 := cc3.at(i, 2, k) + cc3.at(i, 4, k)
			tr4 := cc3.at(i, 4, k) - cc3.at(i, 2, k)
			tr1 := cc3.at(i-1, 1, k) - cc3.at(i-1, 3, k)
			tr2 := cc3.at(i-1, 1, k) + cc3.at(i-1, 3, k)
			ti4 := cc3.at(i-1, 2, k) - cc3.at(i-1, 4, k)
			tr3 := cc3.at(i-1, 2, k) + cc3.at(i-1, 4, k)
			ch3.set(i-1, k, 1, tr2+tr3)
			cr3 := tr2 - tr3
			ch3.set(i, k, 1, ti2+ti3)
			ci3 := ti2 - ti3
			cr2 := tr1 + tr4
			cr4 := tr1 - tr4
			ci2 := ti1 + ti4
			ci4 := ti1 - ti4
			ch3.set(i-1, k, 2, wa1.at(i-1)*cr2-wa1.at(i)*ci2)
			ch3.set(i, k, 2, wa1.at(i-1)*ci2+wa1.at(i)*cr2)
			ch3.set(i-1, k, 3, wa2.at(i-1)*cr3-wa2.at(i)*ci3)
			ch3.set(i, k, 3, wa2.at(i-1)*ci3+wa2.at(i)*cr3)
			ch3.set(i-1, k, 4, wa3.at(i-1)*cr4-wa3.at(i)*ci4)
			ch3.set(i, k, 4, wa3.at(i-1)*ci4+wa3.at(i)*cr4)
		}
	}
}

func passb5(ido, l1 int, cc, ch []float64, wa1, wa2, wa3, wa4 oneArray) {
	const (
		tr11 = 0.309016994374947
		ti11 = 0.951056516295154
		tr12 = -0.809016994374947
		ti12 = 0.587785252292473
	)

	cc3 := newThreeArray(ido, 5, l1, cc)
	ch3 := newThreeArray(ido, l1, 5, ch)

	if ido == 2 {
		for k := 1; k <= l1; k++ {
			ti5 := cc3.at(2, 2, k) - cc3.at(2, 5, k)
			ti2 := cc3.at(2, 2, k) + cc3.at(2, 5, k)
			ti4 := cc3.at(2, 3, k) - cc3.at(2, 4, k)
			ti3 := cc3.at(2, 3, k) + cc3.at(2, 4, k)
			tr5 := cc3.at(1, 2, k) - cc3.at(1, 5, k)
			tr2 := cc3.at(1, 2, k) + cc3.at(1, 5, k)
			tr4 := cc3.at(1, 3, k) - cc3.at(1, 4, k)
			tr3 := cc3.at(1, 3, k) + cc3.at(1, 4, k)
			ch3.set(1, k, 1, cc3.at(1, 1, k)+tr2+tr3)
			ch3.set(2, k, 1, cc3.at(2, 1, k)+ti2+ti3)
			cr2 := cc3.at(1, 1, k) + tr11*tr2 + tr12*tr3
			ci2 := cc3.at(2, 1, k) + tr11*ti2 + tr12*ti3
			cr3 := cc3.at(1, 1, k) + tr12*tr2 + tr11*tr3
			ci3 := cc3.at(2, 1, k) + tr12*ti2 + tr11*ti3
			cr5 := ti11*tr5 + ti12*tr4
			ci5 := ti11*ti5 + ti12*ti4
			cr4 := ti12*tr5 - ti11*tr4
			ci4 := ti12*ti5 - ti11*ti4
			ch3.set(1, k, 2, cr2-ci5)
			ch3.set(1, k, 5, cr2+ci5)
			ch3.set(2, k, 2, ci2+cr5)
			ch3.set(2, k, 3, ci3+cr4)
			ch3.set(1, k, 3, cr3-ci4)
			ch3.set(1, k, 4, cr3+ci4)
			ch3.set(2, k, 4, ci3-cr4)
			ch3.set(2, k, 5, ci2-cr5)
		}
		return
	}
	for k := 1; k <= l1; k++ {
		for i := 2; i <= ido; i += 2 {
			ti5 := cc3.at(i, 2, k) - cc3.at(i, 5, k)
			ti2 := cc3.at(i, 2, k) + cc3.at(i, 5, k)
			ti4 := cc3.at(i, 3, k) - cc3.at(i, 4, k)
			ti3 := cc3.at(i, 3, k) + cc3.at(i, 4, k)
			tr5 := cc3.at(i-1, 2, k) - cc3.at(i-1, 5, k)
			tr2 := cc3.at(i-1, 2, k) + cc3.at(i-1, 5, k)
			tr4 := cc3.at(i-1, 3, k) - cc3.at(i-1, 4, k)
			tr3 := cc3.at(i-1, 3, k) + cc3.at(i-1, 4, k)
			ch3.set(i-1, k, 1, cc3.at(i-1, 1, k)+tr2+tr3)
			ch3.set(i, k, 1, cc3.at(i, 1, k)+ti2+ti3)
			cr2 := cc3.at(i-1, 1, k) + tr11*tr2 + tr12*tr3
			ci2 := cc3.at(i, 1, k) + tr11*ti2 + tr12*ti3
			cr3 := cc3.at(i-1, 1, k) + tr12*tr2 + tr11*tr3
			ci3 := cc3.at(i, 1, k) + tr12*ti2 + tr11*ti3
			cr5 := ti11*tr5 + ti12*tr4
			ci5 := ti11*ti5 + ti12*ti4
			cr4 := ti12*tr5 - ti11*tr4
			ci4 := ti12*ti5 - ti11*ti4
			dr3 := cr3 - ci4
			dr4 := cr3 + ci4
			di3 := ci3 + cr4
			di4 := ci3 - cr4
			dr5 := cr2 + ci5
			dr2 := cr2 - ci5
			di5 := ci2 - cr5
			di2 := ci2 + cr5
			ch3.set(i-1, k, 2, wa1.at(i-1)*dr2-wa1.at(i)*di2)
			ch3.set(i, k, 2, wa1.at(i-1)*di2+wa1.at(i)*dr2)
			ch3.set(i-1, k, 3, wa2.at(i-1)*dr3-wa2.at(i)*di3)
			ch3.set(i, k, 3, wa2.at(i-1)*di3+wa2.at(i)*dr3)
			ch3.set(i-1, k, 4, wa3.at(i-1)*dr4-wa3.at(i)*di4)
			ch3.set(i, k, 4, wa3.at(i-1)*di4+wa3.at(i)*dr4)
			ch3.set(i-1, k, 5, wa4.at(i-1)*dr5-wa4.at(i)*di5)
			ch3.set(i, k, 5, wa4.at(i-1)*di5+wa4.at(i)*dr5)
		}
	}
}

func passb(ido, ip, l1, idl1 int, cc, c1, c2, ch, ch2 []float64, wa oneArray) (nac bool) {
	cc3 := newThreeArray(ido, ip, l1, cc)
	c13 := newThreeArray(ido, l1, ip, c1)
	ch3 := newThreeArray(ido, l1, ip, ch)
	c2m := newTwoArray(idl1, ip, c2)
	ch2m := newTwoArray(idl1, ip, ch2)

	idot := ido / 2
	ipp2 := ip + 2
	ipph := (ip + 1) / 2
	idp := ip * ido

	if ido < l1 {
		for j := 2; j <= ipph; j++ {
			jc := ipp2 - j
			for i := 1; i <= ido; i++ {
				for k := 1; k <= l1; k++ {
					ch3.set(i, k, j, cc3.at(i, j, k)+cc3.at(i, jc, k))
					ch3.set(i, k, jc, cc3.at(i, j, k)-cc3.at(i, jc, k))
				}
			}
		}
		for i := 1; i <= ido; i++ {
			for k := 1; k <= l1; k++ {
				ch3.set(i, k, 1, cc3.at(i, 1, k))
			}
		}
	} else {
		for j := 2; j <= ipph; j++ {
			jc := ipp2 - j
			for k := 1; k <= l1; k++ {
				for i := 1; i <= ido; i++ {
					ch3.set(i, k, j, cc3.at(i, j, k)+cc3.at(i, jc, k))
					ch3.set(i, k, jc, cc3.at(i, j, k)-cc3.at(i, jc, k))
				}
			}
		}
		for k := 1; k <= l1; k++ {
			for i := 1; i <= ido; i++ {
				ch3.set(i, k, 1, cc3.at(i, 1, k))
			}
		}
	}

	idl := 2 - ido
	inc := 0
	for l := 2; l <= ipph; l++ {
		lc := ipp2 - l
		idl += ido
		for ik := 1; ik <= idl1; ik++ {
			c2m.set(ik, l, ch2m.at(ik, 1)+wa.at(idl-1)*ch2m.at(ik, 2))
			c2m.set(ik, lc, wa.at(idl)*ch2m.at(ik, ip))
		}
		idlj := idl
		inc += ido
		for j := 3; j <= ipph; j++ {
			jc := ipp2 - j
			idlj += inc
			if idlj > idp {
				idlj -= idp
			}
			war := wa.at(idlj - 1)
			wai := wa.at(idlj)
			for ik := 1; ik <= idl1; ik++ {
				c2m.set(ik, l, c2m.at(ik, l)+war*ch2m.at(ik, j))
				c2m.set(ik, lc, c2m.at(ik, lc)+wai*ch2m.at(ik, jc))
			}
		}
	}

	for j := 2; j <= ipph; j++ {
		for ik := 1; ik <= idl1; ik++ {
			ch2m.set(ik, 1, ch2m.at(ik, 1)+ch2m.at(ik, j))
		}
	}

	for j := 2; j <= ipph; j++ {
		jc := ipp2 - j
		for ik := 2; ik <= idl1; ik += 2 {
			ch2m.set(ik-1, j, c2m.at(ik-1, j)-c2m.at(ik, jc))
			ch2m.set(ik-1, jc, c2m.at(ik-1, j)+c2m.at(ik, jc))
			ch2m.set(ik, j, c2m.at(ik, j)+c2m.at(ik-1, jc))
			ch2m.set(ik, jc, c2m.at(ik, j)-c2m.at(ik-1, jc))
		}
	}

	if ido == 2 {
		return true
	}

	for ik := 1; ik <= idl1; ik++ {
		c2m.set(ik, 1, ch2m.at(ik, 1))
	}

	for j := 2; j <= ip; j++ {
		for k := 1; k <= l1; k++ {
			c13.set(1, k, j, ch3.at(1, k, j))
			c13.set(2, k, j, ch3.at(2, k, j))
		}
	}

	if idot > l1 {
		idj := 2 - ido
		for j := 2; j <= ip; j++ {
			idj += ido
			for k := 1; k <= l1; k++ {
				idij := idj
				for i := 4; i <= ido; i += 2 {
					idij += 2
					c13.set(i-1, k, j, wa.at(idij-1)*ch3.at(i-1, k, j)-wa.at(idij)*ch3.at(i, k, j))
					c13.set(i, k, j, wa.at(idij-1)*ch3.at(i, k, j)+wa.at(idij)*ch3.at(i-1, k, j))
				}
			}
		}
		return false
	}

	idij := 0
	for j := 2; j <= ip; j++ {
		idij += 2
		for i := 4; i <= ido; i += 2 {
			idij += 2
			for k := 1; k <= l1; k++ {
				c13.set(i-1, k, j, wa.at(idij-1)*ch3.at(i-1, k, j)-wa.at(idij)*ch3.at(i, k, j))
				c13.set(i, k, j, wa.at(idij-1)*ch3.at(i, k, j)+wa.at(idij)*ch3.at(i-1, k, j))
			}
		}
	}
	return false
}