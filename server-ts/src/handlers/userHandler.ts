import supabase from "src/utils/supachaseClient";

/*
 *
 * get user data from supabase
 * function automatically inserts missing information from the "public" user db table
 * and also displays information from the "auth" database from supabase
 *
 */
export const getUser = async (jwt: string) => {
  const {
    data: { user },
  } = await supabase.auth.getUser(jwt);

  //todo: manipulate string to add full name and organisation info from the
  //"public" user database and not only from the "auth" database from supabase

  return user;
};
