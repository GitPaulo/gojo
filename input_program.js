// Variables & Arithmetic
var x = 1;
var y = 1 + 1;

// Conditional Statements
if (x > y) {
    console.log("x greater than y")
} else if (x < y) {
    console.log("x less than or equal to y")
} else {
    console.log("x equal to y")
}

// Switch
var x = 2;
switch (x) {
    case 1: {
        console.log("[switch] x is 1");
        break;
    } case 2: {
        console.log("[switch] x is 2");
        break;
    } default: {
        console.log("[switch] x is something else");
    }
}


// Boolean Operators
console.log(true && false);
console.log(true || false);

// While Loop
let i = 0;
while (i < 5) {
    console.log("while loop!");
    i = i + 1;
}

// Array
var arr = [1, 2, 3, 4, 5];
console.log(arr[2]);

// Built in Functions
console.log("Hello World!");
console.log(Math.sqrt(16));
console.log(Math.pow(3, 2));
