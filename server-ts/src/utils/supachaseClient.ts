import { createClient } from "@supabase/supabase-js";
import dotenv from "dotenv";

//load .env variables
dotenv.config();

const supabaseUrl = process.env.SUPABASE_PROJECT_URL!;
const supabaseAnonKey = process.env.SUPABASE_ANON_PUBLIC_KEY!;
const supabase = createClient(supabaseUrl, supabaseAnonKey);

export default supabase;
