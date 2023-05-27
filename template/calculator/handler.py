from function_execute_platform import status

def handle(input_data, parameters):
    status.update_progress(10)
    first = float(input_data["first"])
    second = float(input_data["second"])

    typ = str(parameters["type"])
    coeff = float(parameters["coeff"])

    status.update_progress(50)
    result = 0.0
    if typ == "sum":
        result = first + second + coeff
    elif typ == "multi":
        return first * second * coeff
    elif typ == "division":
        return first / second / coeff
    
    status.update_progress(100)
    
    return {"result": result}
