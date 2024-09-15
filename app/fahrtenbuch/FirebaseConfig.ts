// Import the functions you need from the SDKs you need

import { initializeApp } from "firebase/app";
import { getAnalytics } from "firebase/analytics";
import { getAuth } from "firebase/auth";

// TODO: Add SDKs for Firebase products that you want to use

// https://firebase.google.com/docs/web/setup#available-libraries

// Your web app's Firebase configuration

// For Firebase JS SDK v7.20.0 and later, measurementId is optional

const firebaseConfig = {
  apiKey: "AIzaSyBuW5DsAk_4KgCwixoYN-MMF8cJtOnHqac",
  authDomain: "fahrtenbuch-57c51.firebaseapp.com",
  projectId: "fahrtenbuch-57c51",
  storageBucket: "fahrtenbuch-57c51.appspot.com",
  messagingSenderId: "125796236587",
  appId: "1:125796236587:web:4373904655d2655212d96e",
  measurementId: "G-WC6ZQFXGM4",
};

// Initialize Firebase

export const firebaseApp = initializeApp(firebaseConfig);
export const firebaseAnalytics = getAnalytics(firebaseApp);
export const firebaseAuth = getAuth(firebaseApp);
