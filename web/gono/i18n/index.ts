"use client";
import i18next from "i18next";
import { initReactI18next } from "react-i18next";
import LanguageDetector from "i18next-browser-languagedetector";
import zh from "./trans/zh.json";
import en from "./trans/en.json";

i18next
  .use(initReactI18next)
  // detect user language
  // learn more: https://github.com/i18next/i18next-browser-languageDetector
  .use(LanguageDetector)
  .init({
    supportedLngs: ["zh", "en"],
    fallbackLng: "en",
    debug: true,
    resources: {
      en: {
        translation: en,
      },
      zh: {
        translation: zh,
      },
    },
    interpolation: {
      escapeValue: false, //React 默认防XSS
    },
  });

export default i18next;
