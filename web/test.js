import './main.js'

const testZip = globalThis.wasmZipModule.newZip()
const testZip2 = globalThis.wasmZipModule.newZip()
const testZip3 = globalThis.wasmZipModule.newZip()
const testZip4 = globalThis.wasmZipModule.newZip()
const testBuffer = new Uint8Array([1,2,3,5,6,7])
const result = globalThis.wasmZipModule.addFile(testZip, 'testFile', testBuffer)
const final = globalThis.wasmZipModule.getFile(testZip)
console.log(final)