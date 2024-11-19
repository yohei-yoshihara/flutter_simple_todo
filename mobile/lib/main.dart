import 'dart:convert';

import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import "package:logging/logging.dart";
import 'package:fluttertoast/fluttertoast.dart';

import 'package:simple_todo/tasks.dart';
import 'package:simple_todo/logging.dart';
import 'package:simple_todo/session.dart';

void main() {
  Logger.root.level = Level.ALL;
  Logger.root.onRecord.listen((LogRecord rec) {
    print('${rec.level.name}: ${rec.time}: ${rec.message}');
  });

  var store = SessionStore();
  store.headers['content-type'] = 'application/json';

  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Todo App',
      theme: ThemeData(
        colorScheme: ColorScheme.fromSeed(seedColor: Colors.green),
        useMaterial3: true,
      ),
      routes: <String, WidgetBuilder>{
        '/login': (BuildContext context) => const LoginPage(),
        '/tasks': (BuildContext context) => const TasksPage(),
      },
      initialRoute: '/login',
    );
  }
}

class LoginPage extends StatefulWidget {
  const LoginPage({super.key});

  @override
  State<LoginPage> createState() => _LoginPageState();
}

class _LoginPageState extends State<LoginPage> {
  final _formKey = GlobalKey<FormState>();
  final _usernameController = TextEditingController();
  final _passwordController = TextEditingController();

  void _submit(BuildContext context) async {
    if (!context.mounted) return;

    if (_formKey.currentState != null && !_formKey.currentState!.validate()) {
      return;
    }

    log.info("login request");
    var url = Uri.http("localhost:8000", "/api/login");

    String body = jsonEncode({
      "username": _usernameController.text,
      "password": _passwordController.text,
    });

    var store = SessionStore();
    http.Response response = await http.post(url, headers: store.headers, body: body);
    store.updateCookie(response);

    log.info('Response status: ${response.statusCode}');
    log.info('Response body: ${response.body}');

    if (response.statusCode != 200 && context.mounted) {
      Fluttertoast.showToast(
        msg: "ログインできませんでした",
        toastLength: Toast.LENGTH_SHORT,
        gravity: ToastGravity.BOTTOM,
        backgroundColor: Colors.red.shade300,
        textColor: Colors.white,
      );
      return;
    }

    if (response.statusCode == 200 && context.mounted) {
      Navigator.of(context).pushReplacementNamed("/tasks");
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        backgroundColor: Theme.of(context).colorScheme.inversePrimary,
        title: const Text("ログイン"),
      ),
      body: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Form(
          key: _formKey,
          child: Column(
            children: [
              TextFormField(
                controller: _usernameController,
                decoration: const InputDecoration(labelText: "ユーザー名"),
                validator: (value) {
                  if (value == null || value.isEmpty) {
                    return "ユーザー名を入力してください";
                  }
                  return null;
                },
              ),
              TextFormField(
                  controller: _passwordController,
                  obscureText: true,
                  decoration: const InputDecoration(labelText: "パスワード"),
                  validator: (value) {
                    if (value == null || value.isEmpty) {
                      return "パスワードを入力してください";
                    }
                    return null;
                  }),
              Padding(
                padding: const EdgeInsets.all(10.0),
                child: ElevatedButton(
                  onPressed: () {
                    _submit(context);
                  },
                  child: const Text('Login'),
                ),
              )
            ],
          ),
        ),
      ),
    );
  }
}
