import supabase from "src/utils/supachaseClient";

/*
 *
 * Start a Ride
 *
 * Takes in
 * 1. jwt => to check for the user
 *
 * 2. carId => assign the ride to a car id
 *
 * 3. ride From => a location string with the street name,
 * postal code, etc... to identy the starting location
 *
 */
interface StartRideResponse {
  status: "success" | "error";
  message: string;
  data?: any;
  error?: string;
}

export const startRide = async (
  jwt: string,
  carId: string,
  rideFrom: string
): Promise<StartRideResponse> => {
  const { data, error } = await supabase.auth.getUser(jwt);
  if (error) {
    return {
      message: "",
      status: "success",
    };
  }

  return {
    message: "",
    status: "success",
  };
};
