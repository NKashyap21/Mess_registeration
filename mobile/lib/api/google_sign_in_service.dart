import 'package:google_sign_in/google_sign_in.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:flutter_dotenv/flutter_dotenv.dart';
import 'api_service.dart';

final GoogleSignIn _googleSignIn = GoogleSignIn(
  scopes: ['email', 'profile'],
  clientId: dotenv.env['GOOGLE_CLIENT_ID'],
);

Future<bool> signInWithGoogle() async {
  try {
    await _googleSignIn.signOut();

    final account = await _googleSignIn.signIn();
    if (account == null) {
      return false;
    }

    final auth = await account.authentication;

    if (auth.idToken == null) {
      throw Exception('ID Token is null! Check your client ID configuration.');
    }

    // Use the simple POST /api/login endpoint for mobile
    final result = await ApiService.loginWithGoogle(auth.idToken!);

    return result;
  } catch (e) {
    try {
      await _googleSignIn.signOut();
    } catch (signOutError) {
      // Ignore sign out errors
    }
    return false;
  }
}

Future<void> signOut() async {
  try {
    // Sign out from Google (local)
    await _googleSignIn.signOut();
  } catch (e) {
    print('Google sign out failed: $e');
    // Continue with other logout steps
  }

  try {
    // Remove local token
    final prefs = await SharedPreferences.getInstance();
    await prefs.remove('jwt_token');
  } catch (e) {
    print('Failed to remove local token: $e');
    // Continue with backend logout
  }

  try {
    // Call backend logout endpoint (with timeout)
    await ApiService.logout();
  } catch (e) {
    print('Backend logout failed: $e');
    // This is not critical - local logout is sufficient
  }
}
