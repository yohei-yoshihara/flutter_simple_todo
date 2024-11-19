import 'package:http/http.dart' as http;

class SessionStore {
  Map<String, String> headers = {};

  static final SessionStore _instance = SessionStore._internal();

  factory SessionStore() {
    return _instance;
  }

  SessionStore._internal();

  void updateCookie(http.Response response) {
    String? rawCookie = response.headers['set-cookie'];
    if (rawCookie != null) {
      int index = rawCookie.indexOf(';');
      headers['cookie'] = (index == -1) ? rawCookie : rawCookie.substring(0, index);
    }
  }
}
