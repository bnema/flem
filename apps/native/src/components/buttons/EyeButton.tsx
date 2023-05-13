import React, { useCallback, useRef, useState } from 'react'
import { StyleSheet, Pressable, Animated } from 'react-native'
import FontAwesome from 'react-native-vector-icons/FontAwesome'
import { COLORS } from '../../utils/Params'

type TypeProps = {
    // onPressIn: (event: GestureResponderEvent) => void,
    onPressIn: any,
    name: string,
    size: number,
    color: string,
}

export default function EyeButton({ onPressIn, name, size, color }: TypeProps) {
    
    const [buttonColor, setButtonColor] = useState(false)
    
    const scale = useRef(new Animated.Value(1)).current

    const AnimatedScale = useCallback((newValue) => {
        Animated.spring(scale, {
            toValue: newValue,
            friction: 4,
            useNativeDriver: true,
        }).start()
    }, [scale])

    return (
        <Pressable onPressIn={() => { onPressIn(); setButtonColor(!buttonColor); AnimatedScale(.8) }} onPressOut={() => AnimatedScale(1)}>
            {!buttonColor
                ? <Animated.View style={[styles.button, { transform: [{ scale }] }]}>
                    <FontAwesome name={name} size={size} color={color} />
                </Animated.View>
                : <Animated.View style={[styles.button, { transform: [{ scale }] }, { borderColor: COLORS.YUP }, { opacity: 1 }]}>
                    <FontAwesome name={name} size={size} color={color} />
                </Animated.View>
            }
        </Pressable>
    )
}

const styles = StyleSheet.create({
    button: {
        justifyContent: 'center',
        alignItems: 'center',
        width: 60,
        height: 60,
        bottom: 80,
        backgroundColor: COLORS.UI,
        borderWidth: 1,
        borderColor: COLORS.BUTTON,
        borderRadius: 40,
        elevation: 5,
        padding: 'auto',
        opacity: .6,
    },
})