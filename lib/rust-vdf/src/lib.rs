#![deny(warnings)]

extern crate vdf;
use vdf::{VDFParams, WesolowskiVDFParams, VDF};

extern crate libc;
use std::ffi::CStr;
use libc::{c_uint, size_t, c_uchar};
use std::{ptr, slice};

#[no_mangle]
pub extern "C" fn hello(name: *const libc::c_char) {
    let buf_name = unsafe { CStr::from_ptr(name).to_bytes() };
    let str_name = String::from_utf8(buf_name.to_vec()).unwrap();
    println!("Hello {}!", str_name);
}

#[no_mangle]
pub extern "C" fn execute(
    difficulty: c_uint,
    input:  *const c_uchar, input_size: size_t,
    output: *mut   c_uchar, output_size: size_t,
    size_in_bits: size_t
) {
    // just checking...
    let bufsize = ((size_in_bits + 7) >> 3) * 2 + 4;  // output and proof; dunno why the extra 4
    assert_eq!(output_size, bufsize);

    let input_slice = unsafe { slice::from_raw_parts(input as *const u8, input_size as usize) };
    let output_slice: &mut[u8] = unsafe { slice::from_raw_parts_mut(output as *mut u8, output_size as usize) };

    let wesolowski_vdf = WesolowskiVDFParams(size_in_bits as u16).new();

    // solve yields tuple of   challenge, difficulty as usize, self.int_size_bits

    let solution = &wesolowski_vdf.solve(input_slice, u64::from(difficulty)).unwrap()[..];

    assert_eq!(output_size, solution.len());

    output_slice.copy_from_slice(solution);
}

#[no_mangle]
pub extern "C" fn verify(
    difficulty: c_uint,
    input:  *const c_uchar,  input_size: size_t,
    proof:  *const c_uchar,  proof_size: size_t,
    size_in_bits: size_t
) -> u8 /*bool*/ {

    let _in = CVec {ptr:input, len:input_size };

    let input_slice = unsafe { slice::from_raw_parts(input as *const u8, input_size as usize) };
    let proof_slice = unsafe { slice::from_raw_parts(proof as *const u8, proof_size as usize) };

    let wesolowski_vdf = WesolowskiVDFParams(size_in_bits as u16).new();
    
    
    let success = wesolowski_vdf.verify(input_slice, difficulty as u64, proof_slice).is_ok();
    
    if success {
        return 1;
    }
    return 0
}

struct CVec {
    ptr: *const u8,
    len: usize,
}

impl std::ops::Deref for CVec {
    type Target = [u8];

    fn deref(&self) -> &[u8] {
        unsafe { slice::from_raw_parts(self.ptr, self.len) }
    }
}

impl Drop for CVec {
    fn drop(&mut self) {
        // unsafe { deallocate_data(self.ptr) };
    }
}

fn _get_vec() -> CVec {
    let   ptr = ptr::null();
    let   len = 0;

    //unsafe {
        // allocate_data(&mut ptr, &mut len);
        assert!(!(ptr as *const u8).is_null());
        assert!(len >= 0);

        CVec {
            ptr,
            len: len as usize,
        }
    //}
}





// extern crate vdf;
// use vdf::{InvalidProof, PietrzakVDFParams, VDFParams, WesolowskiVDFParams, VDF};
// const CORRECT_SOLUTION: &[u8] =
//     b"\x00\x52\x71\xe8\xf9\xab\x2e\xb8\xa2\x90\x6e\x85\x1d\xfc\xb5\x54\x2e\x41\x73\xf0\x16\
//     \xb8\x5e\x29\xd4\x81\xa1\x08\xdc\x82\xed\x3b\x3f\x97\x93\x7b\x7a\xa8\x24\x80\x11\x38\
//     \xd1\x77\x1d\xea\x8d\xae\x2f\x63\x97\xe7\x6a\x80\x61\x3a\xfd\xa3\x0f\x2c\x30\xa3\x4b\
//     \x04\x0b\xaa\xaf\xe7\x6d\x57\x07\xd6\x86\x89\x19\x3e\x5d\x21\x18\x33\xb3\x72\xa6\xa4\
//     \x59\x1a\xbb\x88\xe2\xe7\xf2\xf5\xa5\xec\x81\x8b\x57\x07\xb8\x6b\x8b\x2c\x49\x5c\xa1\
//     \x58\x1c\x17\x91\x68\x50\x9e\x35\x93\xf9\xa1\x68\x79\x62\x0a\x4d\xc4\xe9\x07\xdf\x45\
//     \x2e\x8d\xd0\xff\xc4\xf1\x99\x82\x5f\x54\xec\x70\x47\x2c\xc0\x61\xf2\x2e\xb5\x4c\x48\
//     \xd6\xaa\x5a\xf3\xea\x37\x5a\x39\x2a\xc7\x72\x94\xe2\xd9\x55\xdd\xe1\xd1\x02\xae\x2a\
//     \xce\x49\x42\x93\x49\x2d\x31\xcf\xf2\x19\x44\xa8\xbc\xb4\x60\x89\x93\x06\x5c\x9a\x00\
//     \x29\x2e\x8d\x3f\x46\x04\xe7\x46\x5b\x4e\xee\xfb\x49\x4f\x5b\xea\x10\x2d\xb3\x43\xbb\
//     \x61\xc5\xa1\x5c\x7b\xdf\x28\x82\x06\x88\x5c\x13\x0f\xa1\xf2\xd8\x6b\xf5\xe4\x63\x4f\
//     \xdc\x42\x16\xbc\x16\xef\x7d\xac\x97\x0b\x0e\xe4\x6d\x69\x41\x6f\x9a\x9a\xce\xe6\x51\
//     \xd1\x58\xac\x64\x91\x5b";

// fn main() {
//     let pietrzak_vdf = PietrzakVDFParams(2048).new();
//     assert_eq!(
//         &pietrzak_vdf.solve(b"\xaa", 100).unwrap()[..],
//         CORRECT_SOLUTION
//     );
//     assert!(pietrzak_vdf.verify(b"\xaa", 100, CORRECT_SOLUTION).is_ok());
// }