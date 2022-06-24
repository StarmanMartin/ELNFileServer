<html lang="en">
<head>
    <title>Add new User/Project</title>
    <style>
        .form-container {
            margin: 5px;
            border: 1px solid black;
            border-radius: 5px;
            padding: 10px;
        }
        .msg {
            margin: 15px;
            border-radius: 5px;
            border: 1px solid black;
            padding: 5px;
            text-align: center;
            font-weight: bold;
        }
        .msg.msg-success {
            background-color:aquamarine;
            border-color: chartreuse;
        }
        .msg.msg-err {
            background-color:crimson;
            border-color: orangered;
            color: white;
        }
    </style>
</head>
<body>
<h1>Add new User/Project</h1>
<div class="msg-container">
    {{if .Msg}}
    {{range $val := .Msg}}
    <p class="msg msg-success">{{$val}}</p>
    {{end}}
    {{end}}
    {{if .Err}}
    <p class="msg msg-err">{{.Err}}</p>
    {{end}}
</div>
<div class="form-container">
<form method="post">

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
</div>
</body>
</html>