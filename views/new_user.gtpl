<html>
<head>
    <title></title>
</head>
<body>
<form method="post">
    {{if .Msg}}
    {{range $val := .Msg}}
             <p>{{$val}}</p>
    {{end}}
    {{end}}
    {{if .Err}}
        <p>{{.Err}}</p>
    {{end}}
    <hr>
    <label>Username:
        <input type="text" name="username" value="{{.User}}" required>
    </label>
    <br>
    <br>
    <label>Password:
        <input type="password" name="password" required>
    </label>
    <br>
    <br>
    <label>Project:
        <input type="text" name="project"  value="{{.Project}}" required>
    </label>
    <hr>
    <input type="submit" value="Add Project/User">
</form>
</body>
</html>