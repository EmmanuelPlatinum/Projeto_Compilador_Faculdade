#include <stdio.h>

int main() {
    int y, x;
    x = 10;
    y = 20;
    if (x > y) {
        printf("%d\n", x);
    } else {
        printf("%d\n", y);
    }
    while (x < 100) {
        x = x + 1;
    }
    return 0;
}
