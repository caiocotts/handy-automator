import AsyncStorage from "@react-native-async-storage/async-storage";

const USER_KEY ='handy_automator_user';
const SESSION_KEY ='handy_automator_session';

export type StoredUser = {
    fullName: string;
    email: string;
    password: string;
};

export async function registerUser(user: StoredUser) {
    const existingUser = await AsyncStorage.getItem(USER_KEY);

    if (existingUser) {
        const parsedUser: StoredUser = JSON.parse(existingUser);
        if (parsedUser.email.toLowerCase() === user.email.toLowerCase()) {
            throw new Error('An account with this email already exists');
        }
    }
    await AsyncStorage.setItem(USER_KEY, JSON.stringify(user));
    await AsyncStorage.setItem(SESSION_KEY, JSON.stringify({email: user.email}));
}
export async function signInUser(email: string, password: string) {
    const savedUser = await AsyncStorage.getItem(USER_KEY);

    if(!savedUser) {
        throw new Error('No account found. Please create account');
    }
    const parsedUser: StoredUser = JSON.parse(savedUser);
    if (parsedUser.email.toLowerCase() !== email.toLowerCase() || parsedUser.password !== password) {
        throw new Error('Invalid email or password');
    }
    await AsyncStorage.setItem(USER_KEY, JSON.stringify({email: parsedUser.email}));
    return parsedUser;
}
export async function signOutUser() {
    await AsyncStorage.removeItem(SESSION_KEY);
}
export async function getCurrentSession() {
    const session = await AsyncStorage.getItem(SESSION_KEY);
    return session ? JSON.parse(session) : null;
}