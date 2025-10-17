import 'package:flutter/material.dart';

class AppColors {
  // Primary colors
  static const black = Color(0xFF000000);
  static const darkGrey = Color(0xFF2B2B2B);
  static const midGrey = Color(0xFF555555);
  static const lightGrey = Color(0xFF9E9E9E);
  static const offWhite = Color(0xFFF5F5F5);
  static const white = Color(0xFFFFFFFF);

  // Accent colors
  static const orange = Color(0xFFE4572E);
  static const lightOrange = Color(0xFFFF8906);
  static const yellow = Color(0xFFFFD500);
}

class AppTheme {
  static ThemeData get darkTheme {
    return ThemeData(
      brightness: Brightness.dark,
      scaffoldBackgroundColor: AppColors.black,
      primaryColor: AppColors.orange,
      colorScheme: const ColorScheme.dark(
        primary: AppColors.orange,
        secondary: AppColors.lightOrange,
        background: AppColors.black,
        surface: AppColors.darkGrey,
      ),
      elevatedButtonTheme: ElevatedButtonThemeData(
        style: ElevatedButton.styleFrom(
          backgroundColor: AppColors.darkGrey,
          foregroundColor: AppColors.white,
        ),
      ),
    );
  }

  static ThemeData get lightTheme {
    return ThemeData(
      brightness: Brightness.light,
      scaffoldBackgroundColor: AppColors.white,
      primaryColor: AppColors.orange,
      colorScheme: const ColorScheme.light(
        primary: AppColors.orange,
        secondary: AppColors.lightOrange,
        background: AppColors.white,
        surface: AppColors.offWhite,
      ),
      elevatedButtonTheme: ElevatedButtonThemeData(
        style: ElevatedButton.styleFrom(
          backgroundColor: AppColors.orange,
          foregroundColor: AppColors.white,
        ),
      ),
    );
  }
}
