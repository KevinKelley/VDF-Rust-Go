
void execute(
    unsigned int difficulty,
    char* input,  int input_size,   /*32*/
    char* output, int output_size,  /*516...result+proof*/ 
    int   sizeInBits                /* 2048 */
);

char /*bool*/ verify(
    unsigned int  difficulty,
    char* input,  int input_size,   /*32*/
    char* proof, int output_size,  /*516*/
    int   sizeInBits                /*2048*/
);
