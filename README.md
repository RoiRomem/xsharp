please note that most of this is not yet implemented

# Syntax Guide

This document describes the syntax of the programming language. It covers basic constructs, data types, control structures, functions, and other language-specific features.

---

## 1. Basic Syntax

### 1.1 Comments
- Single-line comment: `// This is a comment`
- Multi-line comment:
  ```
  /*
     This is a multi-line comment
  */
  ```

### 1.2 Identifiers
- Identifiers must start with a letter (a-z, A-Z) or an underscore (`_`).
- Can contain letters, digits (0-9), and underscores.
- Case-sensitive.

Example:
```c
int myVariable = 10;  
```

### 1.3 Keywords
Reserved keywords cannot be used as variable names. Example:
```
if, else, while, for, return, class, public, private, static, virtual, override, new, delete, malloc, free
```

---

## 2. Data Types

### 2.1 Primitive Types
- `int` – Integer values
- `float` – Floating point numbers
- `bool` – Boolean (`true` or `false`)
- `char` – Single character
- `string` – Sequence of characters (pointer to first character in memory)

Example:
```c
int age = 16;
float pi = 3.14;
bool isActive = true;
char initial = 'A';
string name = "John";
```

### 2.2 Static Arrays
```c
int numbers[5] = {1, 2, 3, 4, 5};
```

### 2.3 Dynamic Arrays (Lists)
```c
LIST<int> dynamicList;
dynamicList.Add(10);
dynamicList.Remove(0);
```

---

## 3. Control Structures

### 3.1 Conditional Statements
```c
if (condition) {
    // code block
} else {
    // alternative code block
}
```

### 3.2 Loops
**For loop:**
```c
for (int i = 0; i < 10; i++) {
    // loop body
} 
```

**While loop:**
```c
while (condition) {
    // loop body
}
```

---

## 4. Functions
Functions are defined similarly to C# but with C-like syntax.

Example:
```c
public int add(int a, int b) {
    return a + b;
}
```

---

## 5. Classes and Objects

### 5.1 Class Declaration
```c
public class Person {
    private string name;

    public Person(string n) {
        this->name = n;
    }
    
    public string greet() {
        return "Hello, " + this->name;
    }
}
```

### 5.2 Creating an Object
```c
Person* p = new Person("Alice");
string message = p->greet();
delete p;
```

---

## 6. Subclasses and Inheritance
```c
public class Student : Person {
    private int grade;

    public Student(string n, int g) : Person(n) {
        this->grade = g;
    }
}
```

---

## 7. Memory Management
Memory allocation follows C and C++ conventions, allowing for both static and dynamic memory handling.

### 7.1 Creating Memory on the Heap
```c
int* p = new int;
*p = 42;
```

### 7.2 Freeing Memory
```c
delete p;
```

### 7.3 Allocating and Freeing Objects
```c
Person* p = new Person("John");
delete p;
```

### 7.4 Allocating and Freeing Arrays
```c
int* arr = new int[10];
delete[] arr;
```

### 7.5 Using Raw Pointers
```c
int* ptr = (int*)malloc(sizeof(int));
*ptr = 25;
free(ptr);
```

### 7.6 LIST<> for Dynamic Arrays
```c
LIST<int> numbers;
numbers.Add(10);
numbers.Add(20);
numbers.Remove(0);
```

---

## 8. Input and Output
### 8.1 Printing to Console
```c
printf("Hello, %s!", name);
```

### 8.2 User Input
```c
int age = stdin.int();
string userName = stdin.string();
```

---

This document serves as a quick reference for the language's syntax. Let me know if you need modifications!

