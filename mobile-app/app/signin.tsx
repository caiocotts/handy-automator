import { useState } from 'react';
import { Alert, StyleSheet, TextInput, TouchableOpacity, useColorScheme} from 'react-native';
import { Link, Stack, useRouter } from 'expo-router';
import { ThemedText } from '@/components/themed-text';
import { ThemedView } from '@/components/themed-view';
import {signInUser} from "../services/auth";


export default function SigninScreen() {
    const colorScheme = useColorScheme();
    const router = useRouter();

    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [loading, setLoading] = useState(false);

    const theme = {
        text: colorScheme === 'dark' ? '#fff' : '#000',
        inputBg: colorScheme === 'dark' ? '#111' : '#fff',
        border: colorScheme === 'dark' ? '#555' : '#ccc',
        placeholder: colorScheme === 'dark' ? '#888' : '#666',
        buttonBg: '#0a7ea4',
    };

    const handleSignIn = async () => {
        if (!email.trim || !password.trim) {
            Alert.alert('Missing fields', 'Please enter your email address and password');
            return;
        }
        try {
            setLoading(true);
            await signInUser(email.trim(), password);
            Alert.alert('Success', 'Sign in successfully');
            router.replace('/intro');
        } catch (error) {
            const message = error instanceof Error ? error.message : 'Something went wrong';
            Alert.alert('Sign In Failed', message);
        } finally {
            setLoading(false);
        }
    };
    return (
        <>
            <Stack.Screen options={{ title: 'Sign In',
            headerStyle: {backgroundColor:'#000'},
            headerTintColor: '#fff',
            headerTitleStyle: {fontWeight:'bold'},}}/>
            <ThemedView style={styles.container}>
                <ThemedView style={styles.card}>
                    <ThemedText type="title">Sign in</ThemedText>
                    <ThemedText style={styles.subtitle}>Access your Handy Automator Account</ThemedText>
                    <TextInput placeholder="Email" placeholderTextColor={theme.placeholder} style={[styles.input, {color: theme.text, backgroundColor: theme.inputBg, borderColor: theme.border},]} autoCapitalize="none" keyboardType="email-address" value={email} onChangeText={setEmail}/>
                    <TextInput placeholder="Password" placeholderTextColor={theme.placeholder} style={[styles.input, {color: theme.text, backgroundColor: theme.inputBg, borderColor: theme.border},]} secureTextEntry value={password} onChangeText={setPassword}/>
                    <TouchableOpacity style={[styles.button, {backgroundColor: theme.buttonBg}]} onPress={handleSignIn} disabled={loading}>
                        <ThemedText style={styles.buttonText}>
                            {loading ? 'Sign in...' : 'Signin'}
                        </ThemedText>
                    </TouchableOpacity>
                    <Link href="/create-account">
                        <ThemedText type="defaultSemiBold" style={styles.linkText}>Do not have an account? Create one</ThemedText>
                    </Link>
                </ThemedView>
            </ThemedView>
        </>
    )
}
const styles = StyleSheet.create({
    container: {
        flex: 1,
        justifyContent: 'center',
        padding: 24,
        gap: 16,
    },
    card: {
        gap: 16,
    },
    subtitle: {
        fontSize: 16,
        opacity: 0.8,
        marginBottom: 8,
    },
    input: {
        borderWidth: 1,
        borderColor: '#555',
        borderRadius: 12,
        paddingVertical: 14,
        paddingHorizontal: 14,
        fontSize: 16,
    },
    button: {
        backgroundColor: '#0a7ea4',
        paddingVertical: 14,
        borderRadius: 12,
        alignItems: 'center',
        marginTop: 4,
    },
    buttonText: {
        color: 'white',
        fontWeight: 'bold',
        fontSize: 16,
    },
    linkText: {
        marginTop: 8,
    },
});