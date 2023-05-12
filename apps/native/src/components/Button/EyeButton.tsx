import React from 'react'
import { StyleSheet, View, Pressable } from 'react-native'
import FontAwesome from 'react-native-vector-icons/FontAwesome'

type TypeProps = {
    // onPressIn: (event: GestureResponderEvent) => void,
    onPressIn: any,
    name: string,
    size: number,
    color: string,
}

export default function EyeButton({ onPressIn, name, size, color }: TypeProps) {
    return (
        <Pressable onPressIn={onPressIn}>
            <View style={styles.button}>
                <FontAwesome name={name} size={size} color={color} />
            </View>
        </Pressable>
    )
}

const styles = StyleSheet.create({
    button: {
        justifyContent: 'center',
        alignItems: 'center',
        padding: 'auto',
        bottom: 70,
        opacity: .7,
    },
})