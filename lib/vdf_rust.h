// func New(difficulty int, input [32]byte) *VDF {
// func (vdf *VDF) Execute() {
// func (vdf *VDF) GetOutputChannel() chan [516]byte {
// func (vdf *VDF) Verify(proof [516]byte) bool {
// func (vdf *VDF) IsFinished() bool {
// func (vdf *VDF) GetOutput() [516]byte {


void execute(
    unsigned int difficulty,
    char* input,  int input_size,   /*32*/
    char* output, int output_size,  /*516...result+proof? wait wut?*/ 
    int   sizeInBits                /* 2048 */
);

char /*bool*/ verify(
    unsigned int  difficulty,
    char* input,  int input_size,   /*32*/
    char* output, int output_size,  /*516?*/
    char* proof,  int proof_size,   /*516?*/
    int   sizeInBits                /* 2048 */
);
