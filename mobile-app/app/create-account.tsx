import { useState } from 'react';
import { Alert, StyleSheet, TextInput, TouchableOpacity, useColorScheme} from 'react-native';
import { Link, Stack, useRouter } from 'expo-router';
import { ThemedText } from '@/components/themed-text';
import {ThemedView} from "@/components/themed-view";
import { registerUser } from '../services/auth'

export default function CreateAccountScreen() {
    const colorScheme = useColorScheme();
    const router = useRouter();

    const [fullName, setFullName] = useState('');
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [confirmPassword, setConfirmPassword] = useState('');
    const [loading, setLoading] = useState(false);

    const theme = {
        text: colorScheme === 'dark' ? '#fff' : '#000',
        inputBg: colorScheme === 'dark' ? '#111' : '#fff',
        border: colorScheme === 'dark' ? '#888' : '#ccc',
        placeholder: colorScheme === 'dark' ? '#888' : '#666',
        buttonBg: '#0a7ea4'
    };

    const handleRegister = async () => {
        if (!fullName.trim() || !email.trim() || !password || !confirmPassword) {
            Alert.alert('Missing fields', 'Please fill in all fields');
            return;
        }
        const emailRegex = /\S+@\S+\.\S+/;
        if (!emailRegex.test(email.trim())) {
            Alert.alert('Invalid email','Please enter a valid email address');
            return;
        }
        if (password !== confirmPassword) {
            Alert.alert('Passwords mismatch', 'Passwords do not match');
            return;
        }
        try {
            setLoading(true);
            await registerUser({
                fullName: fullName.trim(),
                email: email.trim(),
                password,
            });

            Alert.alert ('Success', 'Account created successfully');
            router.replace('/signin')
        } catch (error) {
            const message = error instanceof Error ? error.message : 'Something went wrong';
            Alert.alert('Registration Failed', message)
        } finally {
            setLoading(false);
        }
    };
    return (
        <>
            <Stack.Screen options={{
                title: 'Create Account',
                headerStyle: {backgroundColor: '#000'},
                headerTintColor: '#fff',
                headerTitleStyle: {fontWeight: 'bold'},
            }}/>
            <ThemedView style={styles.container}>
                <ThemedView style={styles.card}>
                    <ThemedText type={"title"}>Create Account</ThemedText>
                    <ThemedText style={styles.subtitle}>Set up your Handy Automator Account</ThemedText>
                    <TextInput placeholder="Full Name" placeholderTextColor={theme.placeholder} style={[styles.input, {color: theme.text, backgroundColor: theme.inputBg, borderColor: theme.border},]} value={fullName} onChangeText={setFullName}/>
                    <TextInput placeholder="Email" placeholderTextColor={theme.placeholder} style={[styles.input, {color: theme.text, backgroundColor: theme.inputBg, borderColor: theme.border},]} autoCapitalize="none" keyboardType="email-address" value={email} onChangeText={setEmail}/>
                    <TextInput placeholder="Password" placeholderTextColor={theme.placeholder} style={[styles.input, {color: theme.text, backgroundColor: theme.inputBg, borderColor: theme.border},]} secureTextEntry value={password} onChangeText={setPassword}/>
                    <TextInput placeholder="Confirm Password" placeholderTextColor={theme.placeholder} style={[styles.input, {color: theme.text, backgroundColor: theme.inputBg, borderColor: theme.border},]} secureTextEntry value={confirmPassword} onChangeText={setConfirmPassword}/>

                    <TouchableOpacity style={[styles.button, {backgroundColor: theme.buttonBg}]} onPress={handleRegister} disabled={loading}>
                        <ThemedText style={styles.buttonText}>{loading ? 'Registering...' : 'Register'}</ThemedText>
                    </TouchableOpacity>
                    <Link href= "/signin">
                        <ThemedText type="defaultSemiBold" style={styles.linkText}>Already have an account? Sign in</ThemedText>
                    </Link>
                </ThemedView>
            </ThemedView>
        </>
    );
}
const styles = StyleSheet.create({
    container: {
        flex: 1,
        justifyContent: 'center',
        padding: 24,
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