#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include "sha2.c"

int count_zeros(unsigned char* op) {
    unsigned char v;
    unsigned int z;
    int i;
    for (i = 0; i < SHA256_DIGEST_SIZE; i++) {
        if (op[i] != 0) {
            z = i * 8;
            v = op[i];
            if (v < 16) {
                v <<= 4;
                z += 4;
            }
            if (v < 64) {
                v <<= 2;
                z += 2;
            }
            if (v < 128) {
                z += 1;
            }
            return z;
        }
    }
    return SHA256_DIGEST_SIZE * 8;
}

void print_hex(const unsigned char* c, unsigned int len) {
    int i;
    for (i = 0; i < len; i++) {
        printf("%02x", c[i]);
    }
    printf("\n");
}

// Since the nonce starts from 1, the number of iterations = nonce
int get_nonce(char* msg, unsigned int m_len, unsigned int diff) {
    unsigned int z = 0,
        ip_len = m_len + sizeof(uint32_t);
    uint32_t nonce = 0;

    unsigned char *ip = (unsigned char*) malloc(ip_len);
    unsigned char op[SHA256_DIGEST_SIZE];
    memcpy((void *) ip, (void *) msg, m_len);

    do {
        nonce++;
        UNPACK32(nonce, ip + m_len);
        sha256(ip, ip_len, op);
        z = count_zeros(op);
    } while (z < diff);

    print_hex(op, SHA256_DIGEST_SIZE);
    free(ip);
    return nonce;
}
