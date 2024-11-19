import 'dart:convert';

import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'package:fluttertoast/fluttertoast.dart';

import 'package:simple_todo/logging.dart';
import 'package:simple_todo/session.dart';
import 'package:simple_todo/create_task.dart';
import 'package:simple_todo/edit_task.dart';
import 'package:simple_todo/model.dart';
import 'package:simple_todo/loading.dart';

class TasksPage extends StatefulWidget {
  const TasksPage({super.key});

  @override
  State<TasksPage> createState() => _TasksPageState();
}

class _TasksPageState extends State<TasksPage> {
  List<Task> _tasks = [];
  bool _isLoading = false;

  void _getTasks(BuildContext context) async {
    setState(() {
      _isLoading = true;
    });

    log.info("getTasks");
    var url = Uri.http(
      "localhost:8000",
      "/api/tasks",
    );

    var store = SessionStore();
    http.Response response = await http.get(url, headers: store.headers);
    store.updateCookie(response);

    setState(() {
      _isLoading = false;
    });

    log.info('get: Response status: ${response.statusCode}');
    log.info('get: Response body: ${response.body}');

    if (response.statusCode == 401 && context.mounted) {
      Fluttertoast.showToast(
        msg: "セッションがタイムアウトしました。再度ログインしてください。",
        toastLength: Toast.LENGTH_LONG,
        gravity: ToastGravity.BOTTOM,
        backgroundColor: Colors.red.shade300,
        textColor: Colors.white,
      );
      Navigator.of(context).popUntil((route) => route.isFirst);
      Navigator.pushReplacementNamed(context, "/login");
      return;
    }

    if (response.statusCode == 200 && mounted) {
      String s = utf8.decode(response.bodyBytes);
      List<dynamic> json = jsonDecode(s);
      List<Task> tasks = [];
      for (int i = 0; i < json.length; i++) {
        tasks.add(Task.fromJson(json[i]));
      }
      setState(() {
        _tasks = tasks;
      });
    }
  }

  void _onCreate(BuildContext context) {
    if (context.mounted) {
      Navigator.of(context).push(
        MaterialPageRoute(
          builder: (context) => CreateTaskPage(
            onCreated: (task) {
              setState(() {
                _tasks.insert(0, task);
              });
            },
          ),
        ),
      );
    }
  }

  void _onEdit(BuildContext context, int index) {
    if (context.mounted) {
      Navigator.of(context).push(
        MaterialPageRoute(
          builder: (context) => EditTaskPage(
            task: _tasks[index],
            onUpdated: (task) {
              setState(() {
                var index = _tasks.indexWhere((t) => t.id == task.id);
                if (index >= 0) {
                  _tasks[index] = task;
                }
              });
            },
            onDeleted: (task) {
              setState(() {
                var index = _tasks.indexWhere((t) => t.id == task.id);
                if (index >= 0) {
                  _tasks.removeAt(index);
                }
              });
            },
          ),
        ),
      );
    }
  }

  @override
  void initState() {
    log.info("TasksPage: initState");
    _getTasks(context);
    super.initState();
  }

  @override
  void dispose() {
    log.info("TasksPage: disponse");
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text("タスク一覧"),
        leading: IconButton(
          onPressed: () {
            _getTasks(context);
          },
          icon: const Icon(Icons.sync),
        ),
        actions: [
          IconButton(
            onPressed: () {
              _onCreate(context);
            },
            icon: const Icon(Icons.add),
          ),
        ],
      ),
      body: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Stack(
          children: [
            Column(
              children: <Widget>[
                Expanded(
                  child: ListView.builder(
                    itemCount: _tasks.length,
                    itemBuilder: (context, index) => ListTile(
                      leading: const Icon(Icons.task),
                      title: Text(_tasks[index].title),
                      subtitle: Text(dateFormat.format(_tasks[index].dueDate)),
                      trailing: _tasks[index].completed
                          ? const Icon(
                              Icons.check_circle_outline,
                              color: Colors.green,
                            )
                          : null,
                      isThreeLine: true,
                      onTap: () {
                        _onEdit(context, index);
                      },
                    ),
                  ),
                ),
              ],
            ),
            LoadingWidget(isLoading: _isLoading),
          ],
        ),
      ),
    );
  }
}
