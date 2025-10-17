import 'dart:convert';
import 'package:http/http.dart' as http;
import 'package:shared_preferences/shared_preferences.dart';
import 'package:flutter_dotenv/flutter_dotenv.dart';

class ApiService {
  static final String? baseUrl = dotenv.env['BACKEND_URL'];

  static Future<void> saveToken(String token) async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.setString('jwt_token', token);
  }

  static Future<String?> getToken() async {
    final prefs = await SharedPreferences.getInstance();
    return prefs.getString('jwt_token');
  }

  static Future<void> logout() async {
    final url = Uri.parse('$baseUrl/logout');

    try {
      await http
          .post(url, headers: {'Accept': 'application/json'})
          .timeout(
            const Duration(seconds: 5),
            onTimeout: () {
              return http.Response('{"message": "timeout"}', 408);
            },
          );
    } catch (e) {
      print('Logout API call failed: $e');
    }
  }

  static Future<String?> getGoogleLoginRedirect() async {
    final url = Uri.parse('$baseUrl/login');

    try {
      final response = await http.get(
        url,
        headers: {'Accept': 'application/json'},
      );

      if (response.statusCode == 200) {
        final data = jsonDecode(response.body);
        if (data['data'] != null && data['data']['redirect'] != null) {
          return data['data']['redirect'];
        }
      }
      return null;
    } catch (e) {
      return null;
    }
  }

  static Future<bool> loginWithGoogle(String idToken) async {
    final url = Uri.parse('$baseUrl/login');

    try {
      final body = jsonEncode({'token': idToken});
      final response = await http.post(
        url,
        headers: {
          'Content-Type': 'application/json; charset=utf-8',
          'Accept': 'application/json',
        },
        body: body,
      );
      final data = jsonDecode(response.body);

      if (response.statusCode == 200 &&
          data['data'] != null &&
          data['data']['token'] != null) {
        await saveToken(data['data']['token']);
        return true;
      } else {
        if (data['error'] != null) {}
        return false;
      }
    } catch (e) {
      return false;
    }
  }

  static Future<Map<String, dynamic>> getRegistrationStatus() async {
    final token = await getToken();
    final url = Uri.parse('$baseUrl/students/isRegistrationOpen');

    try {
      final headers = {
        'Content-Type': 'application/json',
        'Accept': 'application/json',
        'Authorization': 'Bearer $token',
      };

      final response = await http.get(url, headers: headers);

      return jsonDecode(response.body);
    } catch (e) {
      return {'error': 'Network error: $e'};
    }
  }

  static Future<Map<String, dynamic>> registerMess({required int mess}) async {
    final token = await getToken();

    final url = Uri.parse('$baseUrl/students/registerMess/$mess');

    try {
      final headers = {
        'Accept': 'application/json',
        'Authorization': 'Bearer $token',
      };

      final response = await http.post(url, headers: headers);

      return jsonDecode(response.body);
    } catch (e) {
      return {'error': 'Network error: $e'};
    }
  }

  static Future<Map<String, dynamic>> registerVegMess() async {
    final token = await getToken();

    final url = Uri.parse('$baseUrl/students/registerVegMess');

    try {
      final headers = {
        'Accept': 'application/json',
        'Authorization': 'Bearer $token',
      };

      final response = await http.post(url, headers: headers);

      return jsonDecode(response.body);
    } catch (e) {
      return {'error': 'Network error: $e'};
    }
  }

  static Future<Map<String, dynamic>> getMess() async {
    final token = await getToken();
    final url = Uri.parse('$baseUrl/students/getMess');

    try {
      final headers = {
        'Content-Type': 'application/json',
        'Accept': 'application/json',
        'Authorization': 'Bearer $token',
      };

      final response = await http.get(url, headers: headers);

      return jsonDecode(response.body);
    } catch (e) {
      return {'error': 'Network error: $e'};
    }
  }

  static Future<Map<String, dynamic>> getMessStats() async {
    final token = await getToken();
    final url = Uri.parse('$baseUrl/students/messStats');

    try {
      final headers = {
        'Content-Type': 'application/json',
        'Accept': 'application/json',
        'Authorization': 'Bearer $token',
      };

      final response = await http.get(url, headers: headers);

      print("MESS STATS RESPONSE: ${response.body}");

      return jsonDecode(response.body);
    } catch (e) {
      return {'error': 'Network error: $e'};
    }
  }

  static Future<Map<String, dynamic>> getMessStatsGrouped() async {
    final token = await getToken();
    final url = Uri.parse('$baseUrl/students/messStatsGrouped');

    try {
      final headers = {
        'Content-Type': 'application/json',
        'Accept': 'application/json',
        'Authorization': 'Bearer $token',
      };

      final response = await http.get(url, headers: headers);

      return jsonDecode(response.body);
    } catch (e) {
      return {'error': 'Network error: $e'};
    }
  }

  static Future<Map<String, dynamic>> getUserInfo() async {
    final token = await getToken();
    final url = Uri.parse('$baseUrl/getUser');

    try {
      final headers = {
        'Content-Type': 'application/json',
        'Accept': 'application/json',
        'Authorization': 'Bearer $token',
      };

      print("Token: $token");

      final response = await http.get(url, headers: headers);

      return jsonDecode(response.body);
    } catch (e) {
      return {'error': 'Network error: $e'};
    }
  }

  static Future<Map<String, dynamic>> createSwapRequest({
    required String type, // "friend" or "public"
    required String password,
  }) async {
    final token = await getToken();
    final url = Uri.parse('$baseUrl/students/createSwap');

    try {
      final headers = {
        'Content-Type': 'application/json; charset=utf-8',
        'Accept': 'application/json',
        'Authorization': 'Bearer $token',
      };

      final body = jsonEncode({'type': type, 'password': password});

      final response = await http.post(url, headers: headers, body: body);

      return jsonDecode(response.body);
    } catch (e) {
      return {'error': 'Network error: $e'};
    }
  }

  static Future<Map<String, dynamic>> deleteSwapRequest() async {
    final token = await getToken();
    final url = Uri.parse('$baseUrl/students/deleteSwap');

    try {
      final headers = {
        'Content-Type': 'application/json; charset=utf-8',
        'Accept': 'application/json',
        'Authorization': 'Bearer $token',
      };

      final response = await http.delete(url, headers: headers);

      return jsonDecode(response.body);
    } catch (e) {
      return {'error': 'Network error: $e'};
    }
  }

  static Future<Map<String, dynamic>> getSwaps() async {
    final token = await getToken();
    final url = Uri.parse('$baseUrl/students/getSwaps');

    try {
      final headers = {
        'Content-Type': 'application/json',
        'Accept': 'application/json',
        'Authorization': 'Bearer $token',
      };

      final response = await http.get(url, headers: headers);

      return jsonDecode(response.body);
    } catch (e) {
      return {'error': 'Network error: $e'};
    }
  }

  static Future<Map<String, dynamic>> acceptSwapRequest({
    required String type,
    required int userId,
  }) async {
    final token = await getToken();
    final url = Uri.parse('$baseUrl/students/acceptSwap');

    try {
      final headers = {
        'Content-Type': 'application/json; charset=utf-8',
        'Accept': 'application/json',
        'Authorization': 'Bearer $token',
      };

      final body = jsonEncode({'type': type, 'user_id': userId});

      final response = await http.post(url, headers: headers, body: body);

      return jsonDecode(response.body);
    } catch (e) {
      return {'error': 'Network error: $e'};
    }
  }

  static Future<Map<String, dynamic>> getUserSwapRequest() async {
    final token = await getToken();
    final url = Uri.parse('$baseUrl/students/getSwapByID');

    try {
      final headers = {
        'Content-Type': 'application/json',
        'Accept': 'application/json',
        'Authorization': 'Bearer $token',
      };

      final response = await http.get(url, headers: headers);

      return jsonDecode(response.body);
    } catch (e) {
      return {'error': 'Network error: $e'};
    }
  }
}
