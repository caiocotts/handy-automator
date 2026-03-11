import { Link } from 'expo-router';
import {StyleSheet, TextInput} from 'react-native';

import { ThemedText } from '@/components/themed-text';
import { ThemedView } from '@/components/themed-view';
import {setName} from "@expo/config-plugins/build/ios/Name";
import { View, Text } from 'react-native';
import {text} from "node:stream/consumers";
import {SafeAreaProvider, SafeAreaView} from "react-native-safe-area-context";
import React from "react";


export default function ModalScreen() {
    const [text, onChangeText] = React.useState('Device Name');
    const [number, onChangeNumber] = React.useState('');
  return (
    <ThemedView style={styles.container}>
      <ThemedText type="title">New Device</ThemedText>


        <SafeAreaProvider>
            <SafeAreaView>

                <TextInput
                    style={[styles.input, { color: '#FFFFFF', backgroundColor: '#8041AB'}]}
                    onChangeText={onChangeText}
                    placeholder="Name of Device"
                />
                <TextInput
                    style={[styles.input, { color: '#FFFFFF', backgroundColor: '#8041AB'}]}
                    onChangeText={onChangeNumber}
                    value={number}
                    placeholder="Type of Device"
                    keyboardType="numeric"
                />

                <button
                style={{height: 50, }}>
                    submit
                </button>
            </SafeAreaView>
        </SafeAreaProvider>



      <Link href="/explore" dismissTo style={styles.link}>
        <ThemedText type="link">Cancel</ThemedText>
      </Link>
    </ThemedView>
  );
}

const styles = StyleSheet.create({
    input: {
        height: 40,
        margin: 12,
        borderWidth: 1,
        padding: 10,
    },
  container: {
    flex: 1,
    alignItems: 'center',
    justifyContent: 'center',
    padding: 20,
  },
  link: {
    marginTop: 15,
    paddingVertical: 15,
  },
});
