import { StyleSheet, View } from "react-native"
import Swiper from "../components/swiper/Swiper"
import { COLORS } from '../utils/Params'

export default function DiscoverScreen() {

    return (
        <View style={styles.container}>
            <Swiper />
        </View>
    )
}

const styles = StyleSheet.create({
    container: {
        flex: 1,
        // flexDirection: "column",
        // justifyContent: 'center',
        // alignItems: 'center',
        backgroundColor:COLORS.BACKGROUND,
    }
})