# Rule Engine with Abstract Syntax Tree (AST)

This Rule Engine is a powerful backend application designed to evaluate complex eligibility rules using an Abstract Syntax Tree (AST). It is built with Golang and provides a flexible way to define, combine, and evaluate rules based on attributes like age, department, salary, and experience.

## Features

- **Rule Creation:** Define rules using a string-based format, which is then parsed into an AST for efficient evaluation.
- **Rule Combination:** Combine multiple rules into a single AST, minimizing redundancy and optimizing performance.
- **Rule Evaluation:** Evaluate the combined rules against JSON data, returning `true` or `false` based on whether the data meets the defined criteria.
- **MongoDB Integration:** Store, retrieve, and manage rules using MongoDB, ensuring persistence and scalability.
- **API Endpoints:** Interact with the rule engine via RESTful API endpoints for creating, combining, evaluating, and deleting rules.

## Sample Rule

Example of a complex rule:
```plaintext
((age > 30 AND department = 'Sales') OR (age < 25 AND department = 'Marketing')) AND (salary > 50000 OR experience > 5)
```

## Technology Stack
- **Backend:**  Golang
- **Database:** MongoDB
- **Frontend:** React (with Tailwind CSS for styling)
- **API Framework:** Gorilla Mux


## Installation

To set up the project locally:

1. **Clone the repository:**

    ```bash
    git clone https://github.com/your-username/rule-engine.git
    ```

2. **Navigate to the project directory:**

    ```bash
    cd rule-engine
    ```

3. **Install the necessary dependencies:**

    ```bash
    go mod tidy
    ```

4. **Run the server:**

    ```bash
    go run main.go
    ```
#### Note: Make sure MongoDB is running on the port 27017.

## Example

### Create a Rule

Send a POST request to `/create_rule` with your rule string.

**Example:**

```bash
curl -X POST http://localhost:8080/create_rule -d '{"rule": "age > 30 AND department = \"Sales\""}'
```

### Combine a Rule

Send a POST request to `/combine_rules` with your rules.

**Example:**

```bash
curl -X POST http://localhost:8080/combine_rules -d '{"rules": ["rule1", "rule2"]}'
```

### Evaluate a Rule for the data

Send a POST request to `/evaluate_rule` with your json data.

**Example:**

```bash
curl -X POST http://localhost:8080/evaluate_rule -d '{"data": {"age": 35, "department": "Sales", "salary": 60000}}'
```
### Delete a Rule
Send a DELETE request to /delete_rule/{rule_id} to remove a stored rule from MongoDB.

```bash
curl -X DELETE http://localhost:8080/delete_rule/1
```
