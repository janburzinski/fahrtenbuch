import { View } from "@/components/Themed";
import { firebaseAuth } from "@/FirebaseConfig";
import React, { useState } from "react";

const Login = () => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState("");
  const auth = firebaseAuth;

  return <View></View>;
};

export default Login;
