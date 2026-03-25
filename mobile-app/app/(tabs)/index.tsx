import { Image } from 'expo-image';
import { StyleSheet, TouchableOpacity } from 'react-native';

import { HelloWave } from '@/components/hello-wave';
import ParallaxScrollView from '@/components/parallax-scroll-view';
import { ThemedText } from '@/components/themed-text';
import { ThemedView } from '@/components/themed-view';
import { Link } from 'expo-router';

export default function HomeScreen() {
  return (
    <ParallaxScrollView
      headerBackgroundColor={{ light: '#A1CEDC', dark: '#1D3D47' }}
      headerImage={
        <Image
          source={require('@/assets/images/partial-react-logo.png')}
          style={styles.reactLogo}
        />
      }>
      <ThemedView style={styles.titleContainer}>
        <ThemedText type="title">Welcome to Handy Automator</ThemedText>
        <HelloWave />
      </ThemedView>
      <ThemedView style={styles.stepContainer}>
          <ThemedText type="subtitle">Step 1: Sign In or Create an Account</ThemedText>
          <ThemedText>
              Please <ThemedText type="defaultSemiBold">sign in or create an account</ThemedText> to verify your identity and access the system.
          </ThemedText>
          <ThemedView style={styles.authButtons}>
      <Link href="/signin" asChild>
          <TouchableOpacity style={styles.signInButton}>
              <ThemedText style={styles.buttonText}>Sign in</ThemedText>
          </TouchableOpacity>
      </Link>

      <Link href="/create-account" asChild>
          <TouchableOpacity style={styles.createAccountButton}>
              <ThemedText style={styles.buttonText}>Create Account</ThemedText>
          </TouchableOpacity>
      </Link>
      </ThemedView>
      </ThemedView>
          {/* <ThemedText type="defaultSemiBold">
            {Platform.select({
              ios: 'cmd + d',
              android: 'cmd + m',
              web: 'F12',
            })}
          </ThemedText>{' '}
          to open developer tools. */}

      <ThemedView style={styles.stepContainer}>
        <Link href="/modal">
          <Link.Trigger>
            <ThemedText type="subtitle">Step 2: Gestures</ThemedText>
          </Link.Trigger>
          <Link.Preview />
          <Link.Menu>
            <Link.MenuAction title="Action" icon="cube" onPress={() => alert('Action pressed')} />
            <Link.MenuAction
              title="Share"
              icon="square.and.arrow.up"
              onPress={() => alert('Share pressed')}
            />
            <Link.Menu title="More" icon="ellipsis">
              <Link.MenuAction
                title="Delete"
                icon="trash"
                destructive
                onPress={() => alert('Delete pressed')}
              />
            </Link.Menu>
          </Link.Menu>
        </Link>

        <ThemedText>
          {`Please walk through our custom gesture creation guide.`}
        </ThemedText>
      </ThemedView>
      <ThemedView style={styles.stepContainer}>
        <ThemedText type="subtitle">Step 3: Start Using the System</ThemedText>
        <ThemedText>
          <ThemedText type="defaultSemiBold">Once your account and gestures are configured, you can control your smart home using facial recognition and custom hand gestures</ThemedText>
        </ThemedText>
      </ThemedView>
    </ParallaxScrollView>
  );
}

const styles = StyleSheet.create({
  titleContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 8,
  },
  stepContainer: {
    gap: 8,
    marginBottom: 8,
  },
  reactLogo: {
    height: 178,
    width: 290,
    bottom: 0,
    left: 0,
    position: 'absolute',
  },
    authButtons: {
      marginTop: 12,
        gap: 12,
    },
    signInButton: {
      backgroundColor: '#0a7ea4',
        paddingVertical: 12,
        paddingHorizontal: 16,
        borderRadius: 10,
        alignItems: 'center',
    },
    createAccountButton: {
      backgroundColor: '#1D3D47',
        paddingVertical: 12,
        paddingHorizontal: 16,
        borderRadius: 10,
        alignItems: 'center',
        borderWidth: 1,
        borderColor: '#0a7ea4',
    },
    buttonText: {
      color: 'white',
        fontWeight: 'bold',
    },
});
